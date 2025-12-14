package config

// ThemePresets contains predefined color schemes.
var ThemePresets = map[string]ColorConfig{
	"default": {
		Light: ColorScheme{
			Background:     "#ffffff",
			Foreground:     "#1f2328",
			Muted:          "#57606a",
			Border:         "#d0d7de",
			Accent:         "#0969da",
			CodeBackground: "#f6f8fa",
		},
		Dark: ColorScheme{
			Background:     "#0d1117",
			Foreground:     "#e6edf3",
			Muted:          "#8b949e",
			Border:         "#30363d",
			Accent:         "#2f81f7",
			CodeBackground: "#161b22",
		},
	},

	"ocean": {
		Light: ColorScheme{
			Background:     "#f8fafc",
			Foreground:     "#0f172a",
			Muted:          "#64748b",
			Border:         "#cbd5e1",
			Accent:         "#0284c7",
			CodeBackground: "#f1f5f9",
		},
		Dark: ColorScheme{
			Background:     "#0f172a",
			Foreground:     "#f1f5f9",
			Muted:          "#94a3b8",
			Border:         "#334155",
			Accent:         "#38bdf8",
			CodeBackground: "#1e293b",
		},
	},

	"forest": {
		Light: ColorScheme{
			Background:     "#fafdf7",
			Foreground:     "#1a2e1a",
			Muted:          "#4a6741",
			Border:         "#c2d4b8",
			Accent:         "#2d6a4f",
			CodeBackground: "#f0f7ec",
		},
		Dark: ColorScheme{
			Background:     "#0f1a0f",
			Foreground:     "#e8f5e0",
			Muted:          "#8fbc8b",
			Border:         "#2d4a2d",
			Accent:         "#52b788",
			CodeBackground: "#1a2e1a",
		},
	},
	"sunset": {
		Light: ColorScheme{
			Background:     "#fffbf5",
			Foreground:     "#2d1f1a",
			Muted:          "#7c6a5d",
			Border:         "#e8d5c4",
			Accent:         "#c2410c",
			CodeBackground: "#fef3e6",
		},
		Dark: ColorScheme{
			Background:     "#1a1210",
			Foreground:     "#fef3e6",
			Muted:          "#b8a090",
			Border:         "#3d2e26",
			Accent:         "#fb923c",
			CodeBackground: "#2d1f1a",
		},
	},
	"midnight": {
		Light: ColorScheme{
			Background:     "#f8f9fc",
			Foreground:     "#1e1e2e",
			Muted:          "#6c6c8a",
			Border:         "#ccd0e0",
			Accent:         "#5b5fc7",
			CodeBackground: "#eef0f8",
		},
		Dark: ColorScheme{
			Background:     "#11111b",
			Foreground:     "#cdd6f4",
			Muted:          "#9399b2",
			Border:         "#313244",
			Accent:         "#89b4fa",
			CodeBackground: "#1e1e2e",
		},
	},
	"lavender": {
		Light: ColorScheme{
			Background:     "#faf8ff",
			Foreground:     "#2e1e3e",
			Muted:          "#7c6a8a",
			Border:         "#e0d4ec",
			Accent:         "#7c3aed",
			CodeBackground: "#f5f0fa",
		},
		Dark: ColorScheme{
			Background:     "#1a1625",
			Foreground:     "#f0e8f8",
			Muted:          "#a090b8",
			Border:         "#3d2e4a",
			Accent:         "#a78bfa",
			CodeBackground: "#2e1e3e",
		},
	},
	"rose": {
		Light: ColorScheme{
			Background:     "#fff5f7",
			Foreground:     "#2e1a1e",
			Muted:          "#8a6c74",
			Border:         "#f0d4dc",
			Accent:         "#e11d48",
			CodeBackground: "#fef0f3",
		},
		Dark: ColorScheme{
			Background:     "#1a1012",
			Foreground:     "#fce8ec",
			Muted:          "#b89098",
			Border:         "#4a2e36",
			Accent:         "#fb7185",
			CodeBackground: "#2e1a1e",
		},
	},
	"nord": {
		Light: ColorScheme{
			Background:     "#eceff4",
			Foreground:     "#2e3440",
			Muted:          "#6b7280",
			Border:         "#d8dee9",
			Accent:         "#5e81ac",
			CodeBackground: "#e5e9f0",
		},
		Dark: ColorScheme{
			Background:     "#2e3440",
			Foreground:     "#eceff4",
			Muted:          "#a3a7b1",
			Border:         "#4c566a",
			Accent:         "#88c0d0",
			CodeBackground: "#3b4252",
		},
	},

	"gruvbox": {
		Light: ColorScheme{
			Background:     "#fbf1c7",
			Foreground:     "#3c3836",
			Muted:          "#7c6f64",
			Border:         "#d5c4a1",
			Accent:         "#d65d0e",
			CodeBackground: "#f2e5bc",
		},
		Dark: ColorScheme{
			Background:     "#282828",
			Foreground:     "#ebdbb2",
			Muted:          "#a89984",
			Border:         "#3c3836",
			Accent:         "#fe8019",
			CodeBackground: "#1d2021",
		},
	},

	"dracula": {
		Light: ColorScheme{
			Background:     "#f8f8f2",
			Foreground:     "#282a36",
			Muted:          "#6272a4",
			Border:         "#dcdcdc",
			Accent:         "#bd93f9",
			CodeBackground: "#eff0eb",
		},
		Dark: ColorScheme{
			Background:     "#282a36",
			Foreground:     "#f8f8f2",
			Muted:          "#6272a4",
			Border:         "#44475a",
			Accent:         "#bd93f9",
			CodeBackground: "#1e1f29",
		},
	},

	"solarized": {
		Light: ColorScheme{
			Background:     "#fdf6e3",
			Foreground:     "#073642",
			Muted:          "#657b83",
			Border:         "#eee8d5",
			Accent:         "#268bd2",
			CodeBackground: "#eee8d5",
		},
		Dark: ColorScheme{
			Background:     "#002b36",
			Foreground:     "#eee8d5",
			Muted:          "#93a1a1",
			Border:         "#073642",
			Accent:         "#268bd2",
			CodeBackground: "#073642",
		},
	},

	"mono": {
		Light: ColorScheme{
			Background:     "#ffffff",
			Foreground:     "#111111",
			Muted:          "#6b7280",
			Border:         "#e5e7eb",
			Accent:         "#000000",
			CodeBackground: "#f9fafb",
		},
		Dark: ColorScheme{
			Background:     "#0a0a0a",
			Foreground:     "#f5f5f5",
			Muted:          "#9ca3af",
			Border:         "#262626",
			Accent:         "#ffffff",
			CodeBackground: "#171717",
		},
	},

	"cyberpunk": {
		Light: ColorScheme{
			Background:     "#fdfcff",
			Foreground:     "#1a1025",
			Muted:          "#6b5b95",
			Border:         "#e0d7ff",
			Accent:         "#ff2ea6",
			CodeBackground: "#f4f0ff",
		},
		Dark: ColorScheme{
			Background:     "#0b0014",
			Foreground:     "#f4e9ff",
			Muted:          "#a78bfa",
			Border:         "#2a104a",
			Accent:         "#ff2ea6",
			CodeBackground: "#1a0033",
		},
	},

	"desert": {
		Light: ColorScheme{
			Background:     "#fff8ed",
			Foreground:     "#3b2f2f",
			Muted:          "#8b7355",
			Border:         "#e6d3b1",
			Accent:         "#c0841d",
			CodeBackground: "#fdf1dc",
		},
		Dark: ColorScheme{
			Background:     "#1f1a14",
			Foreground:     "#fdf1dc",
			Muted:          "#c9b28a",
			Border:         "#3d3326",
			Accent:         "#fbbf24",
			CodeBackground: "#2a2219",
		},
	},

	"ice": {
		Light: ColorScheme{
			Background:     "#f0f9ff",
			Foreground:     "#0c1e2c",
			Muted:          "#5b7c99",
			Border:         "#cfe8f3",
			Accent:         "#0ea5e9",
			CodeBackground: "#e6f6ff",
		},
		Dark: ColorScheme{
			Background:     "#020617",
			Foreground:     "#e0f2fe",
			Muted:          "#7dd3fc",
			Border:         "#082f49",
			Accent:         "#38bdf8",
			CodeBackground: "#031525",
		},
	},

	"coffee": {
		Light: ColorScheme{
			Background:     "#faf6f2",
			Foreground:     "#3a2e2a",
			Muted:          "#7a5e54",
			Border:         "#e4d6cd",
			Accent:         "#7c2d12",
			CodeBackground: "#f3ebe5",
		},
		Dark: ColorScheme{
			Background:     "#1f1713",
			Foreground:     "#f5ede6",
			Muted:          "#b8a29a",
			Border:         "#3b2a23",
			Accent:         "#f97316",
			CodeBackground: "#2a1e18",
		},
	},

	"emerald": {
		Light: ColorScheme{
			Background:     "#f0fdf4",
			Foreground:     "#052e16",
			Muted:          "#5b8a6e",
			Border:         "#bbf7d0",
			Accent:         "#059669",
			CodeBackground: "#dcfce7",
		},
		Dark: ColorScheme{
			Background:     "#022c22",
			Foreground:     "#ecfdf5",
			Muted:          "#6ee7b7",
			Border:         "#064e3b",
			Accent:         "#34d399",
			CodeBackground: "#033a2e",
		},
	},

	"amber": {
		Light: ColorScheme{
			Background:     "#fffbeb",
			Foreground:     "#3b2f0b",
			Muted:          "#8a6d3b",
			Border:         "#fde68a",
			Accent:         "#d97706",
			CodeBackground: "#fef3c7",
		},
		Dark: ColorScheme{
			Background:     "#1c1402",
			Foreground:     "#fef3c7",
			Muted:          "#facc15",
			Border:         "#3f2f05",
			Accent:         "#fbbf24",
			CodeBackground: "#2a1f05",
		},
	},

	"matrix": {
		Light: ColorScheme{
			Background:     "#f6fff8",
			Foreground:     "#042f1a",
			Muted:          "#4d7c5f",
			Border:         "#c6f6d5",
			Accent:         "#16a34a",
			CodeBackground: "#dcfce7",
		},
		Dark: ColorScheme{
			Background:     "#020f07",
			Foreground:     "#d1fae5",
			Muted:          "#4ade80",
			Border:         "#134e2a",
			Accent:         "#22c55e",
			CodeBackground: "#052e16",
		},
	},

	"vscode-dark": {
		Light: ColorScheme{
			Background:     "#ffffff",
			Foreground:     "#1e1e1e",
			Muted:          "#6b7280",
			Border:         "#e5e7eb",
			Accent:         "#007acc",
			CodeBackground: "#f3f3f3",
		},
		Dark: ColorScheme{
			Background:     "#1e1e1e",
			Foreground:     "#d4d4d4",
			Muted:          "#9da1a6",
			Border:         "#2d2d2d",
			Accent:         "#3794ff",
			CodeBackground: "#252526",
		},
	},

	"carbon": {
		Light: ColorScheme{
			Background:     "#f7f7f7",
			Foreground:     "#1a1a1a",
			Muted:          "#5f6368",
			Border:         "#dadce0",
			Accent:         "#111827",
			CodeBackground: "#eeeeee",
		},
		Dark: ColorScheme{
			Background:     "#0f0f0f",
			Foreground:     "#e5e5e5",
			Muted:          "#9ca3af",
			Border:         "#262626",
			Accent:         "#f9fafb",
			CodeBackground: "#1a1a1a",
		},
	},

	"sakura": {
		Light: ColorScheme{
			Background:     "#fff1f2",
			Foreground:     "#3f1d2a",
			Muted:          "#9f5b72",
			Border:         "#fecdd3",
			Accent:         "#be185d",
			CodeBackground: "#ffe4e6",
		},
		Dark: ColorScheme{
			Background:     "#1f0a12",
			Foreground:     "#ffe4e6",
			Muted:          "#fda4af",
			Border:         "#3f0f1f",
			Accent:         "#fb7185",
			CodeBackground: "#2a0d18",
		},
	},

	"terminal": {
		Light: ColorScheme{
			Background:     "#fcfcfc",
			Foreground:     "#000000",
			Muted:          "#6b7280",
			Border:         "#d1d5db",
			Accent:         "#16a34a",
			CodeBackground: "#f4f4f5",
		},
		Dark: ColorScheme{
			Background:     "#000000",
			Foreground:     "#e5e7eb",
			Muted:          "#9ca3af",
			Border:         "#1f2937",
			Accent:         "#22c55e",
			CodeBackground: "#020617",
		},
	},
}
