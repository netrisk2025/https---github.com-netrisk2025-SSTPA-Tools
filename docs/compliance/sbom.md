# SSTPA Tool Software Bill of Materials

This SBOM tracks open source software downloaded or integrated into the SSTPA Tool application during development. Refresh it with `make sbom-generate`; verify it with `make sbom-check`.

The SBOM includes npm workspace dependencies, Go modules, Rust/Tauri Cargo packages, and Docker Compose images. License values are read from package metadata when available and inferred from local license files for Go modules.

Total components: 910
Components requiring license review: 0

| Ecosystem | Software | Version | License | Source |
| --- | --- | --- | --- | --- |
| cargo | adler2 | 2.0.1 | 0BSD OR MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | aho-corasick | 1.1.4 | Unlicense OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | alloc-no-stdlib | 2.0.4 | BSD-3-Clause | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | alloc-stdlib | 0.2.2 | BSD-3-Clause | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | android_system_properties | 0.1.5 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | anyhow | 1.0.102 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | atk | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | atk-sys | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | atomic-waker | 1.1.2 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | autocfg | 1.5.0 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | base64 | 0.21.7 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | base64 | 0.22.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | bit-set | 0.8.0 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | bit-vec | 0.8.0 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | bitflags | 1.3.2 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | bitflags | 2.11.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | block-buffer | 0.10.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | block2 | 0.6.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | brotli | 8.0.2 | BSD-3-Clause AND MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | brotli-decompressor | 5.0.0 | BSD-3-Clause/MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | bumpalo | 3.20.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | bytemuck | 1.25.0 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | byteorder | 1.5.0 | Unlicense OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | bytes | 1.11.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cairo-rs | 0.18.5 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cairo-sys-rs | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | camino | 1.2.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cargo_metadata | 0.19.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cargo_toml | 0.22.3 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cargo-platform | 0.1.9 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cc | 1.2.60 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cesu8 | 1.1.0 | Apache-2.0/MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cfb | 0.7.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cfg-expr | 0.15.8 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cfg-if | 1.0.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | chrono | 0.4.44 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | combine | 4.6.7 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | convert_case | 0.4.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cookie | 0.18.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | core-foundation | 0.10.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | core-foundation-sys | 0.8.7 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | core-graphics | 0.25.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | core-graphics-types | 0.2.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cpufeatures | 0.2.17 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | crc32fast | 1.5.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | crossbeam-channel | 0.5.15 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | crossbeam-utils | 0.8.21 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | crypto-common | 0.1.7 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cssparser | 0.29.6 | MPL-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cssparser | 0.36.0 | MPL-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | cssparser-macros | 0.6.1 | MPL-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ctor | 0.2.9 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | darling | 0.23.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | darling_core | 0.23.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | darling_macro | 0.23.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | deranged | 0.5.8 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | derive_more | 0.99.20 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | derive_more | 2.1.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | derive_more-impl | 2.1.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | digest | 0.10.7 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dirs | 6.0.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dirs-sys | 0.5.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dispatch2 | 0.3.1 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | displaydoc | 0.2.5 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dlopen2 | 0.8.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dlopen2_derive | 0.4.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dom_query | 0.27.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dpi | 0.1.2 | Apache-2.0 AND MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dtoa | 1.0.11 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dtoa-short | 0.3.5 | MPL-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dunce | 1.0.5 | CC0-1.0 OR MIT-0 OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | dyn-clone | 1.0.20 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | embed_plist | 1.2.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | embed-resource | 3.0.8 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | equivalent | 1.0.2 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | erased-serde | 0.4.10 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | fastrand | 2.4.1 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | fdeflate | 0.3.7 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | field-offset | 0.3.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | find-msvc-tools | 0.1.9 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | flate2 | 1.1.9 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | fnv | 1.0.7 | Apache-2.0 / MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | foldhash | 0.1.5 | Zlib | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | foldhash | 0.2.0 | Zlib | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | foreign-types | 0.5.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | foreign-types-macros | 0.2.3 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | foreign-types-shared | 0.3.1 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | form_urlencoded | 1.2.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futf | 0.1.5 | MIT / Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-channel | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-core | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-executor | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-io | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-macro | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-sink | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-task | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | futures-util | 0.3.32 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | fxhash | 0.2.1 | Apache-2.0/MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gdk | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gdk-pixbuf | 0.18.5 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gdk-pixbuf-sys | 0.18.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gdk-sys | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gdkwayland-sys | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gdkx11 | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gdkx11-sys | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | generic-array | 0.14.7 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | getrandom | 0.1.16 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | getrandom | 0.2.17 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | getrandom | 0.3.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | getrandom | 0.4.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gio | 0.18.4 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gio-sys | 0.18.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | glib | 0.18.5 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | glib-macros | 0.18.5 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | glib-sys | 0.18.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | glob | 0.3.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gobject-sys | 0.18.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gtk | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gtk-sys | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | gtk3-macros | 0.18.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | hashbrown | 0.12.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | hashbrown | 0.15.5 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | hashbrown | 0.17.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | heck | 0.4.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | heck | 0.5.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | hex | 0.4.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | html5ever | 0.29.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | html5ever | 0.38.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | http | 1.4.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | http-body | 1.0.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | http-body-util | 0.1.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | httparse | 1.10.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | hyper | 1.9.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | hyper-util | 0.1.20 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | iana-time-zone | 0.1.65 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | iana-time-zone-haiku | 0.1.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ico | 0.5.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | icu_collections | 2.2.0 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | icu_locale_core | 2.2.0 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | icu_normalizer | 2.2.0 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | icu_normalizer_data | 2.2.0 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | icu_properties | 2.2.0 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | icu_properties_data | 2.2.0 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | icu_provider | 2.2.0 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | id-arena | 2.3.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ident_case | 1.0.1 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | idna | 1.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | idna_adapter | 1.2.1 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | indexmap | 1.9.3 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | indexmap | 2.14.0 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | infer | 0.19.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ipnet | 2.12.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | iri-string | 0.7.12 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | itoa | 1.0.18 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | javascriptcore-rs | 1.1.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | javascriptcore-rs-sys | 1.1.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | jni | 0.21.1 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | jni-sys | 0.3.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | jni-sys | 0.4.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | jni-sys-macros | 0.4.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | js-sys | 0.3.95 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | json-patch | 3.0.1 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | jsonptr | 0.6.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | keyboard-types | 0.7.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | kuchikiki | 0.8.8-speedreader | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | leb128fmt | 0.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | libappindicator | 0.9.0 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | libappindicator-sys | 0.9.0 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | libc | 0.2.186 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | libloading | 0.7.4 | ISC | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | libredox | 0.1.16 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | litemap | 0.8.2 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | lock_api | 0.4.14 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | log | 0.4.29 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | mac | 0.1.1 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | markup5ever | 0.14.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | markup5ever | 0.38.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | match_token | 0.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | matches | 0.1.10 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | memchr | 2.8.0 | Unlicense OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | memoffset | 0.9.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | mime | 0.3.17 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | miniz_oxide | 0.8.9 | MIT OR Zlib OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | mio | 1.2.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | muda | 0.17.2 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ndk | 0.9.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ndk-context | 0.1.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ndk-sys | 0.6.0+11769913 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | new_debug_unreachable | 1.0.6 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | nodrop | 0.1.14 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | num_enum | 0.7.6 | BSD-3-Clause OR MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | num_enum_derive | 0.7.6 | BSD-3-Clause OR MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | num-conv | 0.2.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | num-traits | 0.2.19 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2 | 0.6.4 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-app-kit | 0.3.2 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-core-foundation | 0.3.2 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-core-graphics | 0.3.2 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-encode | 4.1.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-exception-helper | 0.1.1 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-foundation | 0.3.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-io-surface | 0.3.2 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-quartz-core | 0.3.2 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-ui-kit | 0.3.2 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | objc2-web-kit | 0.3.2 | Zlib OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | once_cell | 1.21.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | option-ext | 0.2.0 | MPL-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | pango | 0.18.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | pango-sys | 0.18.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | parking_lot | 0.12.5 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | parking_lot_core | 0.9.12 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | percent-encoding | 2.3.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf | 0.8.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf | 0.10.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf | 0.11.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf | 0.13.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_codegen | 0.8.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_codegen | 0.11.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_codegen | 0.13.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_generator | 0.8.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_generator | 0.10.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_generator | 0.11.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_generator | 0.13.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_macros | 0.10.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_macros | 0.11.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_macros | 0.13.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_shared | 0.8.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_shared | 0.10.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_shared | 0.11.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | phf_shared | 0.13.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | pin-project-lite | 0.2.17 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | pkg-config | 0.3.33 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | plist | 1.8.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | png | 0.17.16 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | potential_utf | 0.1.5 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | powerfmt | 0.2.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ppv-lite86 | 0.2.21 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | precomputed-hash | 0.1.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | prettyplease | 0.2.37 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | proc-macro-crate | 1.3.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | proc-macro-crate | 2.0.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | proc-macro-crate | 3.5.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | proc-macro-error | 1.0.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | proc-macro-error-attr | 1.0.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | proc-macro-hack | 0.5.20+deprecated | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | proc-macro2 | 1.0.106 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | quick-xml | 0.38.4 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | quote | 1.0.45 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | r-efi | 5.3.0 | MIT OR Apache-2.0 OR LGPL-2.1-or-later | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | r-efi | 6.0.0 | MIT OR Apache-2.0 OR LGPL-2.1-or-later | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand | 0.7.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand | 0.8.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand_chacha | 0.2.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand_chacha | 0.3.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand_core | 0.5.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand_core | 0.6.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand_hc | 0.2.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rand_pcg | 0.2.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | raw-window-handle | 0.6.2 | MIT OR Apache-2.0 OR Zlib | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | redox_syscall | 0.5.18 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | redox_users | 0.5.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ref-cast | 1.0.25 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | ref-cast-impl | 1.0.25 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | regex | 1.12.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | regex-automata | 0.4.14 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | regex-syntax | 0.8.10 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | reqwest | 0.13.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rustc_version | 0.4.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rustc-hash | 2.1.2 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | rustversion | 1.0.22 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | same-file | 1.0.6 | Unlicense/MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | schemars | 0.8.22 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | schemars | 0.9.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | schemars | 1.2.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | schemars_derive | 0.8.22 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | scopeguard | 1.2.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | selectors | 0.24.0 | MPL-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | selectors | 0.36.1 | MPL-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | semver | 1.0.28 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde | 1.0.228 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_core | 1.0.228 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_derive | 1.0.228 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_derive_internals | 0.29.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_json | 1.0.149 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_repr | 0.1.20 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_spanned | 0.6.9 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_spanned | 1.1.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_with | 3.18.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde_with_macros | 3.18.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serde-untagged | 0.1.9 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serialize-to-javascript | 0.1.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | serialize-to-javascript-impl | 0.1.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | servo_arc | 0.2.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | servo_arc | 0.4.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | sha2 | 0.10.9 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | shlex | 1.3.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | simd-adler32 | 0.3.9 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | siphasher | 0.3.11 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | siphasher | 1.0.2 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | slab | 0.4.12 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | smallvec | 1.15.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | socket2 | 0.6.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | softbuffer | 0.4.8 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | soup3 | 0.5.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | soup3-sys | 0.5.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | stable_deref_trait | 1.2.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | string_cache | 0.8.9 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | string_cache | 0.9.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | string_cache_codegen | 0.5.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | string_cache_codegen | 0.6.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | strsim | 0.11.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | swift-rs | 1.0.7 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | syn | 1.0.109 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | syn | 2.0.117 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | sync_wrapper | 1.0.2 | Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | synstructure | 0.13.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | system-deps | 6.2.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tao | 0.34.8 | Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tao-macros | 0.1.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | target-lexicon | 0.12.16 | Apache-2.0 WITH LLVM-exception | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri | 2.10.3 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri-build | 2.5.6 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri-codegen | 2.5.5 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri-macros | 2.5.5 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri-runtime | 2.10.1 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri-runtime-wry | 2.10.1 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri-utils | 2.8.3 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tauri-winres | 0.3.5 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tendril | 0.4.3 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tendril | 0.5.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | thiserror | 1.0.69 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | thiserror | 2.0.18 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | thiserror-impl | 1.0.69 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | thiserror-impl | 2.0.18 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | time | 0.3.47 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | time-core | 0.1.8 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | time-macros | 0.2.27 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tinystr | 0.8.3 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tokio | 1.52.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tokio-util | 0.7.18 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml | 0.8.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml | 0.9.12+spec-1.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_datetime | 0.6.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_datetime | 0.7.5+spec-1.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_datetime | 1.1.1+spec-1.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_edit | 0.19.15 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_edit | 0.20.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_edit | 0.25.11+spec-1.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_parser | 1.1.2+spec-1.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | toml_writer | 1.1.1+spec-1.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tower | 0.5.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tower-http | 0.6.8 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tower-layer | 0.3.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tower-service | 0.3.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tracing | 0.1.44 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tracing-core | 0.1.36 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | tray-icon | 0.21.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | try-lock | 0.2.5 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | typeid | 1.0.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | typenum | 1.20.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unic-char-property | 0.9.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unic-char-range | 0.9.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unic-common | 0.9.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unic-ucd-ident | 0.9.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unic-ucd-version | 0.9.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unicode-ident | 1.0.24 | (MIT OR Apache-2.0) AND Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unicode-segmentation | 1.13.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | unicode-xid | 0.2.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | url | 2.5.8 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | urlpattern | 0.3.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | utf-8 | 0.7.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | utf8_iter | 1.0.4 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | uuid | 1.23.1 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | version_check | 0.9.5 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | version-compare | 0.2.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | vswhom | 0.1.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | vswhom-sys | 0.1.3 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | walkdir | 2.5.0 | Unlicense/MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | want | 0.3.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasi | 0.9.0+wasi-snapshot-preview1 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasi | 0.11.1+wasi-snapshot-preview1 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasip2 | 1.0.3+wasi-0.2.9 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasip3 | 0.4.0+wasi-0.3.0-rc-2026-01-06 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-bindgen | 0.2.118 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-bindgen-futures | 0.4.68 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-bindgen-macro | 0.2.118 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-bindgen-macro-support | 0.2.118 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-bindgen-shared | 0.2.118 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-encoder | 0.244.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-metadata | 0.244.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasm-streams | 0.5.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wasmparser | 0.244.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | web_atoms | 0.2.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | web-sys | 0.3.95 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | webkit2gtk | 2.0.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | webkit2gtk-sys | 2.0.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | webview2-com | 0.38.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | webview2-com-macros | 0.8.1 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | webview2-com-sys | 0.38.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winapi | 0.3.9 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winapi-i686-pc-windows-gnu | 0.4.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winapi-util | 0.1.11 | Unlicense OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winapi-x86_64-pc-windows-gnu | 0.4.0 | MIT/Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | window-vibrancy | 0.6.0 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows | 0.61.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_aarch64_gnullvm | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_aarch64_gnullvm | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_aarch64_gnullvm | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_aarch64_msvc | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_aarch64_msvc | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_aarch64_msvc | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_gnu | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_gnu | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_gnu | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_gnullvm | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_gnullvm | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_msvc | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_msvc | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_i686_msvc | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_gnu | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_gnu | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_gnu | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_gnullvm | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_gnullvm | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_gnullvm | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_msvc | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_msvc | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows_x86_64_msvc | 0.53.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-collections | 0.2.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-core | 0.61.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-core | 0.62.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-future | 0.2.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-implement | 0.60.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-interface | 0.59.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-link | 0.1.3 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-link | 0.2.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-numerics | 0.2.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-result | 0.3.4 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-result | 0.4.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-strings | 0.4.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-strings | 0.5.1 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-sys | 0.45.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-sys | 0.59.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-sys | 0.60.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-sys | 0.61.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-targets | 0.42.2 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-targets | 0.52.6 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-targets | 0.53.5 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-threading | 0.1.0 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | windows-version | 0.1.7 | MIT OR Apache-2.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winnow | 0.5.40 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winnow | 0.7.15 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winnow | 1.0.2 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | winreg | 0.55.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wit-bindgen | 0.51.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wit-bindgen | 0.57.1 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wit-bindgen-core | 0.51.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wit-bindgen-rust | 0.51.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wit-bindgen-rust-macro | 0.51.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wit-component | 0.244.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wit-parser | 0.244.0 | Apache-2.0 WITH LLVM-exception OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | writeable | 0.6.3 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | wry | 0.54.4 | Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | x11 | 2.21.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | x11-dl | 2.21.0 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | yoke | 0.8.2 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | yoke-derive | 0.8.2 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zerocopy | 0.8.48 | BSD-2-Clause OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zerocopy-derive | 0.8.48 | BSD-2-Clause OR Apache-2.0 OR MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zerofrom | 0.1.7 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zerofrom-derive | 0.1.7 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zerotrie | 0.2.4 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zerovec | 0.11.6 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zerovec-derive | 0.11.3 | Unicode-3.0 | apps/desktop-shell/src-tauri/Cargo.toml |
| cargo | zmij | 1.0.21 | MIT | apps/desktop-shell/src-tauri/Cargo.toml |
| docker | caddy | 2.9 | Apache-2.0 | infra/docker/compose.yaml |
| docker | grafana/grafana | 11.5.2 | AGPL-3.0-only | infra/docker/compose.yaml |
| docker | grafana/tempo | 2.7.1 | AGPL-3.0-only | infra/docker/compose.yaml |
| docker | neo4j | 5.26-community | GPL-3.0-only | infra/docker/compose.yaml |
| docker | otel/opentelemetry-collector | 0.122.1 | Apache-2.0 | infra/docker/compose.yaml |
| docker | prom/prometheus | v3.2.1 | Apache-2.0 | infra/docker/compose.yaml |
| go | cel.dev/expr | v0.16.0 | Apache-2.0 | backend/go.mod |
| go | cloud.google.com/go/compute/metadata | v0.5.0 | Apache-2.0 | backend/go.mod |
| go | dario.cat/mergo | v1.0.2 | BSD-3-Clause | backend/go.mod |
| go | github.com/AdaLogics/go-fuzz-headers | v0.0.0-20240806141605-e8a1dd7889d6 | Apache-2.0 | backend/go.mod |
| go | github.com/Azure/go-ansiterm | v0.0.0-20250102033503-faa5f7b0171c | MIT | backend/go.mod |
| go | github.com/cenkalti/backoff/v4 | v4.3.0 | MIT | backend/go.mod |
| go | github.com/census-instrumentation/opencensus-proto | v0.4.1 | Apache-2.0 | backend/go.mod |
| go | github.com/cespare/xxhash/v2 | v2.3.0 | MIT | backend/go.mod |
| go | github.com/cncf/xds/go | v0.0.0-20240723142845-024c85f92f20 | Apache-2.0 | backend/go.mod |
| go | github.com/containerd/log | v0.1.0 | Apache-2.0 | backend/go.mod |
| go | github.com/containerd/platforms | v0.2.1 | Apache-2.0 | backend/go.mod |
| go | github.com/cpuguy83/dockercfg | v0.3.2 | MIT | backend/go.mod |
| go | github.com/creack/pty | v1.1.18 | MIT | backend/go.mod |
| go | github.com/creack/pty | v1.1.9 | MIT | tools/reference-pipeline/go.mod |
| go | github.com/davecgh/go-spew | v1.1.1 | ISC | backend/go.mod |
| go | github.com/distribution/reference | v0.6.0 | Apache-2.0 | backend/go.mod |
| go | github.com/docker/docker | v28.0.1+incompatible | Apache-2.0 | backend/go.mod |
| go | github.com/docker/go-connections | v0.6.0 | Apache-2.0 | backend/go.mod |
| go | github.com/docker/go-units | v0.5.0 | Apache-2.0 | backend/go.mod |
| go | github.com/ebitengine/purego | v0.10.0 | Apache-2.0 | backend/go.mod |
| go | github.com/envoyproxy/go-control-plane | v0.13.0 | Apache-2.0 | backend/go.mod |
| go | github.com/envoyproxy/protoc-gen-validate | v1.1.0 | Apache-2.0 | backend/go.mod |
| go | github.com/felixge/httpsnoop | v1.0.4 | MIT | backend/go.mod |
| go | github.com/go-chi/chi/v5 | v5.2.1 | MIT | backend/go.mod |
| go | github.com/go-logr/logr | v1.4.3 | Apache-2.0 | backend/go.mod |
| go | github.com/go-logr/stdr | v1.2.2 | Apache-2.0 | backend/go.mod |
| go | github.com/go-ole/go-ole | v1.2.6 | MIT | backend/go.mod |
| go | github.com/gogo/protobuf | v1.3.2 | BSD-3-Clause | backend/go.mod |
| go | github.com/golang/glog | v1.2.2 | Apache-2.0 | backend/go.mod |
| go | github.com/golang/protobuf | v1.5.3 | BSD-3-Clause | backend/go.mod |
| go | github.com/google/go-cmp | v0.7.0 | BSD-3-Clause | backend/go.mod |
| go | github.com/google/uuid | v1.6.0 | BSD-3-Clause | backend/go.mod |
| go | github.com/grpc-ecosystem/grpc-gateway/v2 | v2.16.0 | BSD-3-Clause | backend/go.mod |
| go | github.com/kisielk/errcheck | v1.5.0 | MIT | backend/go.mod |
| go | github.com/kisielk/gotool | v1.0.0 | MIT | backend/go.mod |
| go | github.com/klauspost/compress | v1.18.5 | Apache-2.0 | backend/go.mod |
| go | github.com/kr/pretty | v0.3.1 | MIT | tools/reference-pipeline/go.mod |
| go | github.com/kr/pty | v1.1.1 | MIT | tools/reference-pipeline/go.mod |
| go | github.com/kr/text | v0.2.0 | MIT | tools/reference-pipeline/go.mod |
| go | github.com/lufia/plan9stats | v0.0.0-20211012122336-39d0f177ccd0 | BSD-3-Clause | backend/go.mod |
| go | github.com/magiconair/properties | v1.8.10 | BSD-style | backend/go.mod |
| go | github.com/Microsoft/go-winio | v0.6.2 | MIT | backend/go.mod |
| go | github.com/moby/docker-image-spec | v1.3.1 | Apache-2.0 | backend/go.mod |
| go | github.com/moby/patternmatcher | v0.6.1 | Apache-2.0 | backend/go.mod |
| go | github.com/moby/sys/sequential | v0.6.0 | Apache-2.0 | backend/go.mod |
| go | github.com/moby/sys/user | v0.4.0 | Apache-2.0 | backend/go.mod |
| go | github.com/moby/sys/userns | v0.1.0 | Apache-2.0 | backend/go.mod |
| go | github.com/moby/term | v0.5.2 | Apache-2.0 | backend/go.mod |
| go | github.com/morikuni/aec | v1.0.0 | MIT | backend/go.mod |
| go | github.com/neo4j/neo4j-go-driver/v5 | v5.28.0 | Apache-2.0 | backend/go.mod |
| go | github.com/opencontainers/go-digest | v1.0.0 | Apache-2.0 | backend/go.mod |
| go | github.com/opencontainers/image-spec | v1.1.1 | Apache-2.0 | backend/go.mod |
| go | github.com/pkg/diff | v0.0.0-20210226163009-20ebb0f2a09e | BSD-3-Clause | tools/reference-pipeline/go.mod |
| go | github.com/pkg/errors | v0.9.1 | BSD-style | backend/go.mod |
| go | github.com/planetscale/vtprotobuf | v0.6.1-0.20240319094008-0393e58bdf10 | BSD-3-Clause | backend/go.mod |
| go | github.com/pmezard/go-difflib | v1.0.0 | BSD-style | backend/go.mod |
| go | github.com/power-devops/perfstat | v0.0.0-20240221224432-82ca36839d55 | MIT | backend/go.mod |
| go | github.com/rogpeppe/go-internal | v1.14.1 | BSD-3-Clause | tools/reference-pipeline/go.mod |
| go | github.com/russross/blackfriday | v1.6.0 | BSD-style | backend/go.mod |
| go | github.com/santhosh-tekuri/jsonschema/v5 | v5.3.1 | Apache-2.0 | backend/go.mod |
| go | github.com/shirou/gopsutil/v4 | v4.25.6 | BSD-3-Clause | backend/go.mod |
| go | github.com/sirupsen/logrus | v1.9.4 | MIT | backend/go.mod |
| go | github.com/stretchr/objx | v0.5.2 | MIT | backend/go.mod |
| go | github.com/stretchr/testify | v1.11.1 | MIT | backend/go.mod |
| go | github.com/testcontainers/testcontainers-go | v0.36.0 | MIT | backend/go.mod |
| go | github.com/tklauser/go-sysconf | v0.3.15 | BSD-3-Clause | backend/go.mod |
| go | github.com/tklauser/numcpus | v0.10.0 | Apache-2.0 | backend/go.mod |
| go | github.com/yuin/goldmark | v1.2.1 | MIT | backend/go.mod |
| go | github.com/yusufpapurcu/wmi | v1.2.4 | MIT | backend/go.mod |
| go | go.opentelemetry.io/auto/sdk | v1.2.1 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp | v0.60.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/otel | v1.41.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/otel/exporters/otlp/otlptrace | v1.19.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp | v1.19.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/otel/metric | v1.41.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/otel/sdk | v1.35.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/otel/sdk/metric | v1.35.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/otel/trace | v1.41.0 | Apache-2.0 | backend/go.mod |
| go | go.opentelemetry.io/proto/otlp | v1.0.0 | Apache-2.0 | backend/go.mod |
| go | golang.org/x/crypto | v0.31.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/mod | v0.12.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/mod | v0.21.0 | BSD-3-Clause | tools/reference-pipeline/go.mod |
| go | golang.org/x/net | v0.33.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/oauth2 | v0.22.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/sync | v0.8.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/sys | v0.31.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/sys | v0.26.0 | BSD-3-Clause | tools/reference-pipeline/go.mod |
| go | golang.org/x/term | v0.27.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/text | v0.21.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/time | v0.0.0-20220210224613-90d013bbcef8 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/tools | v0.13.0 | BSD-3-Clause | backend/go.mod |
| go | golang.org/x/tools | v0.26.0 | BSD-3-Clause | tools/reference-pipeline/go.mod |
| go | golang.org/x/xerrors | v0.0.0-20200804184101-5ec99f83aff1 | BSD-3-Clause | backend/go.mod |
| go | google.golang.org/genproto/googleapis/api | v0.0.0-20240814211410-ddb44dafa142 | Apache-2.0 | backend/go.mod |
| go | google.golang.org/genproto/googleapis/rpc | v0.0.0-20240903143218-8af14fe29dc1 | Apache-2.0 | backend/go.mod |
| go | google.golang.org/grpc | v1.67.0 | Apache-2.0 | backend/go.mod |
| go | google.golang.org/protobuf | v1.34.2 | BSD-3-Clause | backend/go.mod |
| go | gopkg.in/check.v1 | v1.0.0-20201130134442-10cb98267c6c | BSD-style | tools/reference-pipeline/go.mod |
| go | gopkg.in/yaml.v3 | v3.0.1 | Apache-2.0 | tools/reference-pipeline/go.mod |
| go | gotest.tools/v3 | v3.5.2 | Apache-2.0 | backend/go.mod |
| npm | @alloc/quick-lru | 5.2.0 | MIT | package-lock.json |
| npm | @babel/code-frame | 7.29.0 | MIT | package-lock.json |
| npm | @babel/compat-data | 7.29.0 | MIT | package-lock.json |
| npm | @babel/core | 7.29.0 | MIT | package-lock.json |
| npm | @babel/generator | 7.29.1 | MIT | package-lock.json |
| npm | @babel/helper-compilation-targets | 7.28.6 | MIT | package-lock.json |
| npm | @babel/helper-globals | 7.28.0 | MIT | package-lock.json |
| npm | @babel/helper-module-imports | 7.28.6 | MIT | package-lock.json |
| npm | @babel/helper-module-transforms | 7.28.6 | MIT | package-lock.json |
| npm | @babel/helper-plugin-utils | 7.28.6 | MIT | package-lock.json |
| npm | @babel/helper-string-parser | 7.27.1 | MIT | package-lock.json |
| npm | @babel/helper-validator-identifier | 7.28.5 | MIT | package-lock.json |
| npm | @babel/helper-validator-option | 7.27.1 | MIT | package-lock.json |
| npm | @babel/helpers | 7.29.2 | MIT | package-lock.json |
| npm | @babel/parser | 7.29.2 | MIT | package-lock.json |
| npm | @babel/plugin-transform-react-jsx-self | 7.27.1 | MIT | package-lock.json |
| npm | @babel/plugin-transform-react-jsx-source | 7.27.1 | MIT | package-lock.json |
| npm | @babel/template | 7.28.6 | MIT | package-lock.json |
| npm | @babel/traverse | 7.29.0 | MIT | package-lock.json |
| npm | @babel/types | 7.29.0 | MIT | package-lock.json |
| npm | @esbuild/aix-ppc64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/android-arm | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/android-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/android-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/darwin-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/darwin-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/freebsd-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/freebsd-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-arm | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-ia32 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-loong64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-mips64el | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-ppc64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-riscv64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-s390x | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/linux-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/netbsd-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/netbsd-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/openbsd-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/openbsd-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/openharmony-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/sunos-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/win32-arm64 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/win32-ia32 | 0.25.12 | MIT | package-lock.json |
| npm | @esbuild/win32-x64 | 0.25.12 | MIT | package-lock.json |
| npm | @eslint-community/eslint-utils | 4.9.1 | MIT | package-lock.json |
| npm | @eslint-community/regexpp | 4.12.2 | MIT | package-lock.json |
| npm | @eslint/config-array | 0.21.2 | Apache-2.0 | package-lock.json |
| npm | @eslint/config-helpers | 0.4.2 | Apache-2.0 | package-lock.json |
| npm | @eslint/core | 0.17.0 | Apache-2.0 | package-lock.json |
| npm | @eslint/eslintrc | 3.3.5 | MIT | package-lock.json |
| npm | @eslint/js | 9.39.4 | MIT | package-lock.json |
| npm | @eslint/object-schema | 2.1.7 | Apache-2.0 | package-lock.json |
| npm | @eslint/plugin-kit | 0.4.1 | Apache-2.0 | package-lock.json |
| npm | @humanfs/core | 0.19.2 | Apache-2.0 | package-lock.json |
| npm | @humanfs/node | 0.16.8 | Apache-2.0 | package-lock.json |
| npm | @humanfs/types | 0.15.0 | Apache-2.0 | package-lock.json |
| npm | @humanwhocodes/module-importer | 1.0.1 | Apache-2.0 | package-lock.json |
| npm | @humanwhocodes/retry | 0.4.3 | Apache-2.0 | package-lock.json |
| npm | @jridgewell/gen-mapping | 0.3.13 | MIT | package-lock.json |
| npm | @jridgewell/remapping | 2.3.5 | MIT | package-lock.json |
| npm | @jridgewell/resolve-uri | 3.1.2 | MIT | package-lock.json |
| npm | @jridgewell/sourcemap-codec | 1.5.5 | MIT | package-lock.json |
| npm | @jridgewell/trace-mapping | 0.3.31 | MIT | package-lock.json |
| npm | @nodelib/fs.scandir | 2.1.5 | MIT | package-lock.json |
| npm | @nodelib/fs.stat | 2.0.5 | MIT | package-lock.json |
| npm | @nodelib/fs.walk | 1.2.8 | MIT | package-lock.json |
| npm | @rolldown/pluginutils | 1.0.0-beta.27 | MIT | package-lock.json |
| npm | @rollup/rollup-android-arm-eabi | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-android-arm64 | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-darwin-arm64 | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-darwin-x64 | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-freebsd-arm64 | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-freebsd-x64 | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-arm-gnueabihf | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-arm-musleabihf | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-arm64-gnu | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-arm64-musl | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-loong64-gnu | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-loong64-musl | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-ppc64-gnu | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-ppc64-musl | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-riscv64-gnu | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-riscv64-musl | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-s390x-gnu | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-x64-gnu | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-linux-x64-musl | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-openbsd-x64 | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-openharmony-arm64 | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-win32-arm64-msvc | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-win32-ia32-msvc | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-win32-x64-gnu | 4.60.2 | MIT | package-lock.json |
| npm | @rollup/rollup-win32-x64-msvc | 4.60.2 | MIT | package-lock.json |
| npm | @tanstack/query-core | 5.99.2 | MIT | package-lock.json |
| npm | @tanstack/react-query | 5.99.2 | MIT | package-lock.json |
| npm | @tanstack/react-virtual | 3.13.24 | MIT | package-lock.json |
| npm | @tanstack/virtual-core | 3.14.0 | MIT | package-lock.json |
| npm | @types/babel__core | 7.20.5 | MIT | package-lock.json |
| npm | @types/babel__generator | 7.27.0 | MIT | package-lock.json |
| npm | @types/babel__template | 7.4.4 | MIT | package-lock.json |
| npm | @types/babel__traverse | 7.28.0 | MIT | package-lock.json |
| npm | @types/chai | 5.2.3 | MIT | package-lock.json |
| npm | @types/deep-eql | 4.0.2 | MIT | package-lock.json |
| npm | @types/estree | 1.0.8 | MIT | package-lock.json |
| npm | @types/json-schema | 7.0.15 | MIT | package-lock.json |
| npm | @types/node | 22.19.17 | MIT | package-lock.json |
| npm | @types/react | 19.2.14 | MIT | package-lock.json |
| npm | @types/react-dom | 19.2.3 | MIT | package-lock.json |
| npm | @typescript-eslint/eslint-plugin | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/parser | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/project-service | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/scope-manager | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/tsconfig-utils | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/type-utils | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/types | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/typescript-estree | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/utils | 8.59.0 | MIT | package-lock.json |
| npm | @typescript-eslint/visitor-keys | 8.59.0 | MIT | package-lock.json |
| npm | @vitejs/plugin-react | 4.7.0 | MIT | package-lock.json |
| npm | @vitest/expect | 3.2.4 | MIT | package-lock.json |
| npm | @vitest/mocker | 3.2.4 | MIT | package-lock.json |
| npm | @vitest/pretty-format | 3.2.4 | MIT | package-lock.json |
| npm | @vitest/runner | 3.2.4 | MIT | package-lock.json |
| npm | @vitest/snapshot | 3.2.4 | MIT | package-lock.json |
| npm | @vitest/spy | 3.2.4 | MIT | package-lock.json |
| npm | @vitest/utils | 3.2.4 | MIT | package-lock.json |
| npm | acorn | 8.16.0 | MIT | package-lock.json |
| npm | acorn-jsx | 5.3.2 | MIT | package-lock.json |
| npm | ajv | 6.14.0 | MIT | package-lock.json |
| npm | ansi-styles | 4.3.0 | MIT | package-lock.json |
| npm | any-promise | 1.3.0 | MIT | package-lock.json |
| npm | anymatch | 3.1.3 | ISC | package-lock.json |
| npm | arg | 5.0.2 | MIT | package-lock.json |
| npm | argparse | 2.0.1 | Python-2.0 | package-lock.json |
| npm | assertion-error | 2.0.1 | MIT | package-lock.json |
| npm | autoprefixer | 10.5.0 | MIT | package-lock.json |
| npm | balanced-match | 4.0.4 | MIT | package-lock.json |
| npm | balanced-match | 1.0.2 | MIT | package-lock.json |
| npm | baseline-browser-mapping | 2.10.20 | Apache-2.0 | package-lock.json |
| npm | binary-extensions | 2.3.0 | MIT | package-lock.json |
| npm | brace-expansion | 5.0.5 | MIT | package-lock.json |
| npm | brace-expansion | 1.1.14 | MIT | package-lock.json |
| npm | braces | 3.0.3 | MIT | package-lock.json |
| npm | browserslist | 4.28.2 | MIT | package-lock.json |
| npm | cac | 6.7.14 | MIT | package-lock.json |
| npm | callsites | 3.1.0 | MIT | package-lock.json |
| npm | camelcase-css | 2.0.1 | MIT | package-lock.json |
| npm | caniuse-lite | 1.0.30001788 | CC-BY-4.0 | package-lock.json |
| npm | chai | 5.3.3 | MIT | package-lock.json |
| npm | chalk | 4.1.2 | MIT | package-lock.json |
| npm | check-error | 2.1.3 | MIT | package-lock.json |
| npm | chokidar | 3.6.0 | MIT | package-lock.json |
| npm | color-convert | 2.0.1 | MIT | package-lock.json |
| npm | color-name | 1.1.4 | MIT | package-lock.json |
| npm | commander | 4.1.1 | MIT | package-lock.json |
| npm | concat-map | 0.0.1 | MIT | package-lock.json |
| npm | convert-source-map | 2.0.0 | MIT | package-lock.json |
| npm | cross-spawn | 7.0.6 | MIT | package-lock.json |
| npm | cssesc | 3.0.0 | MIT | package-lock.json |
| npm | csstype | 3.2.3 | MIT | package-lock.json |
| npm | cytoscape | 3.33.2 | MIT | package-lock.json |
| npm | debug | 4.4.3 | MIT | package-lock.json |
| npm | deep-eql | 5.0.2 | MIT | package-lock.json |
| npm | deep-is | 0.1.4 | MIT | package-lock.json |
| npm | didyoumean | 1.2.2 | Apache-2.0 | package-lock.json |
| npm | dlv | 1.1.3 | MIT | package-lock.json |
| npm | electron-to-chromium | 1.5.340 | ISC | package-lock.json |
| npm | es-errors | 1.3.0 | MIT | package-lock.json |
| npm | es-module-lexer | 1.7.0 | MIT | package-lock.json |
| npm | esbuild | 0.25.12 | MIT | package-lock.json |
| npm | escalade | 3.2.0 | MIT | package-lock.json |
| npm | escape-string-regexp | 4.0.0 | MIT | package-lock.json |
| npm | eslint | 9.39.4 | MIT | package-lock.json |
| npm | eslint-scope | 8.4.0 | BSD-2-Clause | package-lock.json |
| npm | eslint-visitor-keys | 3.4.3 | Apache-2.0 | package-lock.json |
| npm | eslint-visitor-keys | 5.0.1 | Apache-2.0 | package-lock.json |
| npm | eslint-visitor-keys | 4.2.1 | Apache-2.0 | package-lock.json |
| npm | espree | 10.4.0 | BSD-2-Clause | package-lock.json |
| npm | esquery | 1.7.0 | BSD-3-Clause | package-lock.json |
| npm | esrecurse | 4.3.0 | BSD-2-Clause | package-lock.json |
| npm | estraverse | 5.3.0 | BSD-2-Clause | package-lock.json |
| npm | estree-walker | 3.0.3 | MIT | package-lock.json |
| npm | esutils | 2.0.3 | BSD-2-Clause | package-lock.json |
| npm | expect-type | 1.3.0 | Apache-2.0 | package-lock.json |
| npm | fast-deep-equal | 3.1.3 | MIT | package-lock.json |
| npm | fast-glob | 3.3.3 | MIT | package-lock.json |
| npm | fast-json-stable-stringify | 2.1.0 | MIT | package-lock.json |
| npm | fast-levenshtein | 2.0.6 | MIT | package-lock.json |
| npm | fastq | 1.20.1 | ISC | package-lock.json |
| npm | fdir | 6.5.0 | MIT | package-lock.json |
| npm | file-entry-cache | 8.0.0 | MIT | package-lock.json |
| npm | fill-range | 7.1.1 | MIT | package-lock.json |
| npm | find-up | 5.0.0 | MIT | package-lock.json |
| npm | flat-cache | 4.0.1 | MIT | package-lock.json |
| npm | flatted | 3.4.2 | ISC | package-lock.json |
| npm | fraction.js | 5.3.4 | MIT | package-lock.json |
| npm | framer-motion | 12.38.0 | MIT | package-lock.json |
| npm | fsevents | 2.3.3 | MIT | package-lock.json |
| npm | function-bind | 1.1.2 | MIT | package-lock.json |
| npm | gensync | 1.0.0-beta.2 | MIT | package-lock.json |
| npm | glob-parent | 5.1.2 | ISC | package-lock.json |
| npm | glob-parent | 6.0.2 | ISC | package-lock.json |
| npm | globals | 14.0.0 | MIT | package-lock.json |
| npm | has-flag | 4.0.0 | MIT | package-lock.json |
| npm | hasown | 2.0.3 | MIT | package-lock.json |
| npm | ignore | 7.0.5 | MIT | package-lock.json |
| npm | ignore | 5.3.2 | MIT | package-lock.json |
| npm | import-fresh | 3.3.1 | MIT | package-lock.json |
| npm | imurmurhash | 0.1.4 | MIT | package-lock.json |
| npm | is-binary-path | 2.1.0 | MIT | package-lock.json |
| npm | is-core-module | 2.16.1 | MIT | package-lock.json |
| npm | is-extglob | 2.1.1 | MIT | package-lock.json |
| npm | is-glob | 4.0.3 | MIT | package-lock.json |
| npm | is-number | 7.0.0 | MIT | package-lock.json |
| npm | isexe | 2.0.0 | ISC | package-lock.json |
| npm | jiti | 1.21.7 | MIT | package-lock.json |
| npm | js-tokens | 4.0.0 | MIT | package-lock.json |
| npm | js-tokens | 9.0.1 | MIT | package-lock.json |
| npm | js-yaml | 4.1.1 | MIT | package-lock.json |
| npm | jsesc | 3.1.0 | MIT | package-lock.json |
| npm | json-buffer | 3.0.1 | MIT | package-lock.json |
| npm | json-schema-traverse | 0.4.1 | MIT | package-lock.json |
| npm | json-stable-stringify-without-jsonify | 1.0.1 | MIT | package-lock.json |
| npm | json5 | 2.2.3 | MIT | package-lock.json |
| npm | keyv | 4.5.4 | MIT | package-lock.json |
| npm | levn | 0.4.1 | MIT | package-lock.json |
| npm | lilconfig | 3.1.3 | MIT | package-lock.json |
| npm | lines-and-columns | 1.2.4 | MIT | package-lock.json |
| npm | locate-path | 6.0.0 | MIT | package-lock.json |
| npm | lodash.merge | 4.6.2 | MIT | package-lock.json |
| npm | loose-envify | 1.4.0 | MIT | package-lock.json |
| npm | loupe | 3.2.1 | MIT | package-lock.json |
| npm | lru-cache | 5.1.1 | ISC | package-lock.json |
| npm | magic-string | 0.30.21 | MIT | package-lock.json |
| npm | merge2 | 1.4.1 | MIT | package-lock.json |
| npm | micromatch | 4.0.8 | MIT | package-lock.json |
| npm | minimatch | 10.2.5 | BlueOak-1.0.0 | package-lock.json |
| npm | minimatch | 3.1.5 | ISC | package-lock.json |
| npm | motion-dom | 12.38.0 | MIT | package-lock.json |
| npm | motion-utils | 12.36.0 | MIT | package-lock.json |
| npm | ms | 2.1.3 | MIT | package-lock.json |
| npm | mz | 2.7.0 | MIT | package-lock.json |
| npm | nanoid | 3.3.11 | MIT | package-lock.json |
| npm | natural-compare | 1.4.0 | MIT | package-lock.json |
| npm | node-releases | 2.0.37 | MIT | package-lock.json |
| npm | normalize-path | 3.0.0 | MIT | package-lock.json |
| npm | object-assign | 4.1.1 | MIT | package-lock.json |
| npm | object-hash | 3.0.0 | MIT | package-lock.json |
| npm | optionator | 0.9.4 | MIT | package-lock.json |
| npm | p-limit | 3.1.0 | MIT | package-lock.json |
| npm | p-locate | 5.0.0 | MIT | package-lock.json |
| npm | parent-module | 1.0.1 | MIT | package-lock.json |
| npm | path-exists | 4.0.0 | MIT | package-lock.json |
| npm | path-key | 3.1.1 | MIT | package-lock.json |
| npm | path-parse | 1.0.7 | MIT | package-lock.json |
| npm | pathe | 2.0.3 | MIT | package-lock.json |
| npm | pathval | 2.0.1 | MIT | package-lock.json |
| npm | picocolors | 1.1.1 | ISC | package-lock.json |
| npm | picomatch | 2.3.2 | MIT | package-lock.json |
| npm | picomatch | 4.0.4 | MIT | package-lock.json |
| npm | pify | 2.3.0 | MIT | package-lock.json |
| npm | pirates | 4.0.7 | MIT | package-lock.json |
| npm | postcss | 8.5.10 | MIT | package-lock.json |
| npm | postcss-import | 15.1.0 | MIT | package-lock.json |
| npm | postcss-js | 4.1.0 | MIT | package-lock.json |
| npm | postcss-load-config | 6.0.1 | MIT | package-lock.json |
| npm | postcss-nested | 6.2.0 | MIT | package-lock.json |
| npm | postcss-selector-parser | 6.1.2 | MIT | package-lock.json |
| npm | postcss-value-parser | 4.2.0 | MIT | package-lock.json |
| npm | prelude-ls | 1.2.1 | MIT | package-lock.json |
| npm | prop-types | 15.8.1 | MIT | package-lock.json |
| npm | punycode | 2.3.1 | MIT | package-lock.json |
| npm | queue-microtask | 1.2.3 | MIT | package-lock.json |
| npm | react | 19.2.5 | MIT | package-lock.json |
| npm | react-cytoscapejs | 2.0.0 | MIT | package-lock.json |
| npm | react-dom | 19.2.5 | MIT | package-lock.json |
| npm | react-is | 16.13.1 | MIT | package-lock.json |
| npm | react-refresh | 0.17.0 | MIT | package-lock.json |
| npm | read-cache | 1.0.0 | MIT | package-lock.json |
| npm | readdirp | 3.6.0 | MIT | package-lock.json |
| npm | resolve | 1.22.12 | MIT | package-lock.json |
| npm | resolve-from | 4.0.0 | MIT | package-lock.json |
| npm | reusify | 1.1.0 | MIT | package-lock.json |
| npm | rollup | 4.60.2 | MIT | package-lock.json |
| npm | run-parallel | 1.2.0 | MIT | package-lock.json |
| npm | scheduler | 0.27.0 | MIT | package-lock.json |
| npm | semver | 6.3.1 | ISC | package-lock.json |
| npm | semver | 7.7.4 | ISC | package-lock.json |
| npm | shebang-command | 2.0.0 | MIT | package-lock.json |
| npm | shebang-regex | 3.0.0 | MIT | package-lock.json |
| npm | siginfo | 2.0.0 | ISC | package-lock.json |
| npm | source-map-js | 1.2.1 | BSD-3-Clause | package-lock.json |
| npm | stackback | 0.0.2 | MIT | package-lock.json |
| npm | std-env | 3.10.0 | MIT | package-lock.json |
| npm | strip-json-comments | 3.1.1 | MIT | package-lock.json |
| npm | strip-literal | 3.1.0 | MIT | package-lock.json |
| npm | sucrase | 3.35.1 | MIT | package-lock.json |
| npm | supports-color | 7.2.0 | MIT | package-lock.json |
| npm | supports-preserve-symlinks-flag | 1.0.0 | MIT | package-lock.json |
| npm | tailwindcss | 3.4.19 | MIT | package-lock.json |
| npm | thenify | 3.3.1 | MIT | package-lock.json |
| npm | thenify-all | 1.6.0 | MIT | package-lock.json |
| npm | tinybench | 2.9.0 | MIT | package-lock.json |
| npm | tinyexec | 0.3.2 | MIT | package-lock.json |
| npm | tinyglobby | 0.2.16 | MIT | package-lock.json |
| npm | tinypool | 1.1.1 | MIT | package-lock.json |
| npm | tinyrainbow | 2.0.0 | MIT | package-lock.json |
| npm | tinyspy | 4.0.4 | MIT | package-lock.json |
| npm | to-regex-range | 5.0.1 | MIT | package-lock.json |
| npm | ts-api-utils | 2.5.0 | MIT | package-lock.json |
| npm | ts-interface-checker | 0.1.13 | Apache-2.0 | package-lock.json |
| npm | tslib | 2.8.1 | 0BSD | package-lock.json |
| npm | type-check | 0.4.0 | MIT | package-lock.json |
| npm | typescript | 5.9.3 | Apache-2.0 | package-lock.json |
| npm | typescript-eslint | 8.59.0 | MIT | package-lock.json |
| npm | undici-types | 6.21.0 | MIT | package-lock.json |
| npm | update-browserslist-db | 1.2.3 | MIT | package-lock.json |
| npm | uri-js | 4.4.1 | BSD-2-Clause | package-lock.json |
| npm | util-deprecate | 1.0.2 | MIT | package-lock.json |
| npm | vite | 6.4.2 | MIT | package-lock.json |
| npm | vite-node | 3.2.4 | MIT | package-lock.json |
| npm | vitest | 3.2.4 | MIT | package-lock.json |
| npm | which | 2.0.2 | ISC | package-lock.json |
| npm | why-is-node-running | 2.3.0 | MIT | package-lock.json |
| npm | word-wrap | 1.2.5 | MIT | package-lock.json |
| npm | yallist | 3.1.1 | ISC | package-lock.json |
| npm | yocto-queue | 0.1.0 | MIT | package-lock.json |
| npm | zustand | 5.0.12 | MIT | package-lock.json |
