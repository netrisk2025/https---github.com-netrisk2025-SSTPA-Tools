#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.

fn main() {
  tauri::Builder::default()
    .run(tauri::generate_context!())
    .expect("failed to run tauri application")
}
