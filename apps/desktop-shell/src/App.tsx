import { toolRegistry } from "./app/toolRegistry"

import { shellBackgroundClassName, glassPanelClassName } from "@sstpa/ui"

export default function App() {
  return (
    <main className={`${shellBackgroundClassName} px-6 py-10`}>
      <div className="mx-auto flex max-w-6xl flex-col gap-8">
        <header className={`${glassPanelClassName} p-8`}>
          <p className="text-xs uppercase tracking-[0.35em] text-cyan-300/80">SSTPA Tools</p>
          <h1 className="mt-4 text-4xl font-semibold text-white">Desktop Shell Scaffold</h1>
          <p className="mt-4 max-w-3xl text-sm leading-7 text-slate-300">
            This shell hosts independently developed add-on tools and keeps the data workflow explicit for SoI work,
            reference-data ingestion, and future verification automation.
          </p>
        </header>

        <section className="grid gap-5 md:grid-cols-3">
          {toolRegistry.map((tool) => (
            <article key={tool.manifest.id} className={`${glassPanelClassName} flex flex-col gap-4 p-6`}>
              <div>
                <p className="text-xs uppercase tracking-[0.3em] text-cyan-200/80">{tool.manifest.runtime}</p>
                <h2 className="mt-3 text-xl font-medium text-white">{tool.manifest.name}</h2>
              </div>

              <p className="text-sm leading-6 text-slate-300">{tool.manifest.summary}</p>

              <div className="space-y-3 text-xs text-slate-400">
                <div>
                  <p className="font-semibold uppercase tracking-[0.25em] text-slate-300">Inputs</p>
                  <p>{tool.manifest.inputContracts.join(", ")}</p>
                </div>

                <div>
                  <p className="font-semibold uppercase tracking-[0.25em] text-slate-300">Outputs</p>
                  <p>{tool.manifest.outputContracts.join(", ")}</p>
                </div>
              </div>
            </article>
          ))}
        </section>
      </div>
    </main>
  )
}
