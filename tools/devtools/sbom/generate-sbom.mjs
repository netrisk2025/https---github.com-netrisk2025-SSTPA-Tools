// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import { execFileSync } from "node:child_process"
import { existsSync, mkdirSync, readFileSync, readdirSync, statSync, writeFileSync } from "node:fs"
import { dirname, join, relative, resolve } from "node:path"
import { fileURLToPath } from "node:url"

const repoRoot = resolve(dirname(fileURLToPath(import.meta.url)), "../../..")
const sbomPath = join(repoRoot, "docs/compliance/sbom.md")
const checkMode = process.argv.includes("--check")

const entries = new Map()

function addEntry(entry) {
  if (!entry.name || !entry.version) {
    return
  }

  const normalized = {
    ecosystem: entry.ecosystem,
    name: entry.name,
    version: String(entry.version),
    license: normalizeLicense(entry.license),
    source: entry.source,
  }
  const key = `${normalized.ecosystem}\u0000${normalized.name}\u0000${normalized.version}`
  entries.set(key, normalized)
}

function normalizeLicense(license) {
  if (!license) {
    return "UNKNOWN - REVIEW REQUIRED"
  }

  if (typeof license === "object") {
    return normalizeLicense(license.type || license.name || license.license)
  }

  return String(license).replaceAll("|", "\\|").trim() || "UNKNOWN - REVIEW REQUIRED"
}

function readJSON(path) {
  return JSON.parse(readFileSync(path, "utf8"))
}

function collectNPM() {
  const lockPath = join(repoRoot, "package-lock.json")
  if (!existsSync(lockPath)) {
    return
  }

  const lock = readJSON(lockPath)
  for (const [packagePath, info] of Object.entries(lock.packages || {})) {
    if (!packagePath.startsWith("node_modules/") || !info.version) {
      continue
    }

    const packageJSONPath = join(repoRoot, packagePath, "package.json")
    const packageJSON = existsSync(packageJSONPath) ? readJSON(packageJSONPath) : {}
    const name = packageJSON.name || packagePath.replace(/^node_modules\//, "")
    if (name.startsWith("@sstpa/")) {
      continue
    }

    addEntry({
      ecosystem: "npm",
      name,
      version: packageJSON.version || info.version,
      license: packageJSON.license || info.license,
      source: "package-lock.json",
    })
  }
}

function collectGoModule(modulePath) {
  const cwd = join(repoRoot, modulePath)
  if (!existsSync(join(cwd, "go.mod"))) {
    return
  }

  const output = execFileSync("go", ["list", "-m", "-f", "{{.Path}}\t{{.Version}}\t{{.Dir}}", "all"], {
    cwd,
    env: { ...process.env, GOWORK: "off" },
    encoding: "utf8",
  })

  for (const line of output.trim().split("\n")) {
    const [name, version, dir] = line.split("\t")
    if (!version || version === "(devel)") {
      continue
    }

    const licenseDir = dir || downloadGoModule(cwd, name, version)
    addEntry({
      ecosystem: "go",
      name,
      version,
      license: inferLicenseFromDirectory(licenseDir) || goLicenseOverrides.get(name),
      source: join(modulePath, "go.mod"),
    })
  }
}

function downloadGoModule(cwd, name, version) {
  try {
    const output = execFileSync("go", ["mod", "download", "-json", `${name}@${version}`], {
      cwd,
      env: { ...process.env, GOWORK: "off" },
      encoding: "utf8",
      maxBuffer: 8 * 1024 * 1024,
    })
    return JSON.parse(output).Dir
  } catch {
    return ""
  }
}

function collectGo() {
  collectGoModule("backend")
  collectGoModule("tools/reference-pipeline")
  collectGoModule("tools/devtools/copyright")
}

function collectCargo() {
  const manifestPath = join(repoRoot, "apps/desktop-shell/src-tauri/Cargo.toml")
  if (!existsSync(manifestPath)) {
    return
  }

  try {
    const output = execFileSync(
      "cargo",
      ["metadata", "--format-version", "1", "--manifest-path", manifestPath],
      { cwd: repoRoot, encoding: "utf8", maxBuffer: 32 * 1024 * 1024 },
    )
    const metadata = JSON.parse(output)
    for (const pkg of metadata.packages || []) {
      if (pkg.name === "sstpa-desktop-shell") {
        continue
      }

      addEntry({
        ecosystem: "cargo",
        name: pkg.name,
        version: pkg.version,
        license: pkg.license || inferLicenseFromManifestPath(pkg.manifest_path),
        source: "apps/desktop-shell/src-tauri/Cargo.toml",
      })
    }
  } catch (error) {
    addCargoFallback(manifestPath)
  }
}

function addCargoFallback(manifestPath) {
  const content = readFileSync(manifestPath, "utf8")
  const known = new Map([
    ["serde", "MIT OR Apache-2.0"],
    ["serde_json", "MIT OR Apache-2.0"],
    ["tauri", "Apache-2.0 OR MIT"],
    ["tauri-build", "Apache-2.0 OR MIT"],
  ])

  for (const line of content.split("\n")) {
    const match = line.match(/^([A-Za-z0-9_-]+)\s*=\s*(?:"([^"]+)"|\{\s*version\s*=\s*"([^"]+)")/)
    if (!match) {
      continue
    }

    addEntry({
      ecosystem: "cargo",
      name: match[1],
      version: match[2] || match[3],
      license: known.get(match[1]),
      source: "apps/desktop-shell/src-tauri/Cargo.toml",
    })
  }
}

function collectDockerImages() {
  const composePath = join(repoRoot, "infra/docker/compose.yaml")
  if (!existsSync(composePath)) {
    return
  }

  const licenses = new Map([
    ["caddy", "Apache-2.0"],
    ["neo4j", "GPL-3.0-only"],
    ["otel/opentelemetry-collector", "Apache-2.0"],
    ["prom/prometheus", "Apache-2.0"],
    ["grafana/tempo", "AGPL-3.0-only"],
    ["grafana/grafana", "AGPL-3.0-only"],
  ])

  const content = readFileSync(composePath, "utf8")
  for (const match of content.matchAll(/^\s*image:\s*["']?([^"'\s]+)["']?\s*$/gm)) {
    const image = match[1]
    const [name, version = "latest"] = splitImage(image)
    addEntry({
      ecosystem: "docker",
      name,
      version,
      license: licenses.get(name),
      source: "infra/docker/compose.yaml",
    })
  }
}

const goLicenseOverrides = new Map([
  ["github.com/containerd/typeurl/v2", "Apache-2.0"],
  ["github.com/creack/pty", "MIT"],
  ["github.com/felixge/httpsnoop", "MIT"],
  ["github.com/gogo/protobuf", "BSD-3-Clause"],
  ["github.com/kr/pretty", "MIT"],
  ["github.com/kr/text", "MIT"],
  ["github.com/moby/sys/mount", "Apache-2.0"],
  ["github.com/moby/sys/mountinfo", "Apache-2.0"],
  ["github.com/moby/sys/reexec", "Apache-2.0"],
  ["github.com/russross/blackfriday", "BSD-2-Clause"],
  ["github.com/santhosh-tekuri/jsonschema/v5", "MIT"],
  ["golang.org/x/mod", "BSD-3-Clause"],
  ["golang.org/x/net", "BSD-3-Clause"],
  ["golang.org/x/text", "BSD-3-Clause"],
  ["golang.org/x/time", "BSD-3-Clause"],
  ["golang.org/x/tools", "BSD-3-Clause"],
  ["golang.org/x/xerrors", "BSD-3-Clause"],
  ["google.golang.org/genproto/googleapis/rpc", "Apache-2.0"],
  ["google.golang.org/grpc", "Apache-2.0"],
  ["google.golang.org/protobuf", "BSD-3-Clause"],
])

function splitImage(image) {
  const atDigest = image.split("@")[0]
  const lastSlash = atDigest.lastIndexOf("/")
  const lastColon = atDigest.lastIndexOf(":")
  if (lastColon > lastSlash) {
    return [atDigest.slice(0, lastColon), atDigest.slice(lastColon + 1)]
  }

  return [atDigest, "latest"]
}

function inferLicenseFromManifestPath(manifestPath) {
  if (!manifestPath) {
    return undefined
  }

  return inferLicenseFromDirectory(dirname(manifestPath))
}

function inferLicenseFromDirectory(dir) {
  if (!dir || !existsSync(dir)) {
    return undefined
  }

  const candidates = findLicenseFiles(dir)
  for (const candidate of candidates) {
    const text = readFileSync(candidate, "utf8")
    const inferred = inferLicenseFromText(text)
    if (inferred) {
      return inferred
    }
  }

  return undefined
}

function findLicenseFiles(dir) {
  const direct = readdirSync(dir)
    .filter((name) => /^(license|licence|copying|notice)(\..*)?$/i.test(name))
    .map((name) => join(dir, name))
    .filter((path) => statSync(path).isFile())

  return direct.slice(0, 4)
}

function inferLicenseFromText(text) {
  const normalized = text.toLowerCase()
  const squashed = normalized.replace(/\s+/g, " ")

  if (normalized.includes("apache license") && normalized.includes("version 2.0")) {
    return "Apache-2.0"
  }
  if (normalized.includes("mit license")) {
    return "MIT"
  }
  if (
    squashed.includes("permission is hereby granted, free of charge") &&
    squashed.includes("the software without restriction")
  ) {
    return "MIT"
  }
  if (normalized.includes("isc license")) {
    return "ISC"
  }
  if (normalized.includes("mozilla public license version 2.0")) {
    return "MPL-2.0"
  }
  if (normalized.includes("gnu affero general public license") && normalized.includes("version 3")) {
    return "AGPL-3.0-only"
  }
  if (normalized.includes("gnu general public license") && normalized.includes("version 3")) {
    return "GPL-3.0-only"
  }
  if (normalized.includes("redistribution and use in source and binary forms")) {
    if (normalized.includes("neither the name") || normalized.includes("nor the names")) {
      return "BSD-3-Clause"
    }
    return "BSD-style"
  }
  if (normalized.includes("unlicense")) {
    return "Unlicense"
  }

  return undefined
}

function render() {
  const rows = Array.from(entries.values()).sort((a, b) =>
    `${a.ecosystem} ${a.name}`.localeCompare(`${b.ecosystem} ${b.name}`),
  )

  const unknownCount = rows.filter((row) => row.license.includes("UNKNOWN")).length
  const lines = [
    "# SSTPA Tool Software Bill of Materials",
    "",
    "This SBOM tracks open source software downloaded or integrated into the SSTPA Tool application during development. Refresh it with `make sbom-generate`; verify it with `make sbom-check`.",
    "",
    "The SBOM includes npm workspace dependencies, Go modules, Rust/Tauri Cargo packages, and Docker Compose images. License values are read from package metadata when available and inferred from local license files for Go modules.",
    "",
    `Total components: ${rows.length}`,
    `Components requiring license review: ${unknownCount}`,
    "",
    "| Ecosystem | Software | Version | License | Source |",
    "| --- | --- | --- | --- | --- |",
  ]

  for (const row of rows) {
    lines.push(
      `| ${escapeCell(row.ecosystem)} | ${escapeCell(row.name)} | ${escapeCell(row.version)} | ${escapeCell(row.license)} | ${escapeCell(row.source)} |`,
    )
  }

  lines.push("")
  return lines.join("\n")
}

function escapeCell(value) {
  return String(value).replaceAll("|", "\\|")
}

collectNPM()
collectGo()
collectCargo()
collectDockerImages()

const next = render()
if (checkMode) {
  const current = existsSync(sbomPath) ? readFileSync(sbomPath, "utf8") : ""
  if (current !== next) {
    console.error("SBOM is not current. Run `make sbom-generate`.")
    process.exit(1)
  }
} else {
  mkdirSync(dirname(sbomPath), { recursive: true })
  writeFileSync(sbomPath, next)
  console.log(`wrote ${relative(repoRoot, sbomPath)}`)
}
