# Add-on Tool Boundaries

Add-on tools are expected to evolve more independently than the host shell.

Design rules for this repository:

- Each add-on lives in its own workspace under `addons/`.
- The shell imports add-on manifests, not internal implementation details.
- Shared contracts belong in `packages/addon-sdk/`.
- Shared domain types belong in `packages/domain/`.
- Shared styling primitives belong in `packages/ui/`.
- Add-ons can open as pop-up workflows or panel workflows, but editing authority still follows the active product rules and backend validation.
- Tool-specific persistence, query keys, and UI state should stay local unless another package has a clear reuse need.

This structure lets the Navigator, Requirements, Message Center, and future MBSE or reporting tools be built and tested with less coupling to the main shell.
