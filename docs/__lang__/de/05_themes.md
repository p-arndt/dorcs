---
title: "Themes"
description: "Durchsuchen Sie alle verfügbaren Themes und passen Sie Ihre Dokumentationsseite an."
tags: [themes, customization, styling]
date: 2025-12-13
draft: false
---

# Themes

dorcs kommt mit einer Vielzahl von integrierten Themes. Jedes Theme enthält sowohl helle als auch dunkle Farbschemata mit passender Syntax-Hervorhebung.

## Themes verwenden

**Über Konfigurationsdatei:**
```yaml
theme:
  preset: midnight
  mode: auto  # light, dark oder auto
```

**Über Befehlszeile:**
```bash
./dorcs --theme midnight --theme-mode dark
```

## Verfügbare Themes

## Default

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #ffffff; color: #1f2328; padding: 20px; border: 2px solid #d0d7de; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #0969da; margin-top: 0;">Default Theme (Hell)</h3>
        <p style="color: #57606a;">Ein sauberes, professionelles Theme mit blauer Akzentfarbe. Perfekt für Dokumentation und technische Inhalte.</p>
        <div style="background-color: #f6f8fa; border: 1px solid #d0d7de; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #0969da;">const</span> <span style="color: #1f2328;">theme</span> = <span style="color: #57606a;">"default"</span>;
        </div>
    </div>
    <div style="background-color: #0d1117; color: #e6edf3; padding: 20px; border: 2px solid #30363d; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #2f81f7; margin-top: 0;">Default Theme (Dunkel)</h3>
        <p style="color: #8b949e;">Ein sauberes, professionelles Theme mit blauer Akzentfarbe. Perfekt für Dokumentation und technische Inhalte.</p>
        <div style="background-color: #161b22; border: 1px solid #30363d; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #2f81f7;">const</span> <span style="color: #e6edf3;">theme</span> = <span style="color: #8b949e;">"default"</span>;
        </div>
    </div>
</div>

## Ocean

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #f8fafc; color: #0f172a; padding: 20px; border: 2px solid #cbd5e1; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #0284c7; margin-top: 0;">Ocean Theme (Hell)</h3>
        <p style="color: #64748b;">Inspiriert vom tiefblauen Meer, dieses Theme zeichnet sich durch kühle Blautöne und ruhige Farben aus.</p>
        <div style="background-color: #f1f5f9; border: 1px solid #cbd5e1; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #0284c7;">const</span> <span style="color: #0f172a;">theme</span> = <span style="color: #64748b;">"ocean"</span>;
        </div>
    </div>
    <div style="background-color: #0f172a; color: #f1f5f9; padding: 20px; border: 2px solid #334155; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #38bdf8; margin-top: 0;">Ocean Theme (Dunkel)</h3>
        <p style="color: #94a3b8;">Inspiriert vom tiefblauen Meer, dieses Theme zeichnet sich durch kühle Blautöne und ruhige Farben aus.</p>
        <div style="background-color: #1e293b; border: 1px solid #334155; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #38bdf8;">const</span> <span style="color: #f1f5f9;">theme</span> = <span style="color: #94a3b8;">"ocean"</span>;
        </div>
    </div>
</div>

## Forest

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fafdf7; color: #1a2e1a; padding: 20px; border: 2px solid #c2d4b8; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #2d6a4f; margin-top: 0;">Forest Theme (Hell)</h3>
        <p style="color: #4a6741;">Von der Natur inspirierte Grüntöne, die Ihrer Dokumentation ein frisches, organisches Gefühl verleihen.</p>
        <div style="background-color: #f0f7ec; border: 1px solid #c2d4b8; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #2d6a4f;">const</span> <span style="color: #1a2e1a;">theme</span> = <span style="color: #4a6741;">"forest"</span>;
        </div>
    </div>
    <div style="background-color: #0f1a0f; color: #e8f5e0; padding: 20px; border: 2px solid #2d4a2d; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #52b788; margin-top: 0;">Forest Theme (Dunkel)</h3>
        <p style="color: #8fbc8b;">Von der Natur inspirierte Grüntöne, die Ihrer Dokumentation ein frisches, organisches Gefühl verleihen.</p>
        <div style="background-color: #1a2e1a; border: 1px solid #2d4a2d; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #52b788;">const</span> <span style="color: #e8f5e0;">theme</span> = <span style="color: #8fbc8b;">"forest"</span>;
        </div>
    </div>
</div>

## Sunset

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fffbf5; color: #2d1f1a; padding: 20px; border: 2px solid #e8d5c4; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #c2410c; margin-top: 0;">Sunset Theme (Hell)</h3>
        <p style="color: #7c6a5d;">Warme Orangen- und Erdtöne, die die goldene Stunde evozieren.</p>
        <div style="background-color: #fef3e6; border: 1px solid #e8d5c4; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #c2410c;">const</span> <span style="color: #2d1f1a;">theme</span> = <span style="color: #7c6a5d;">"sunset"</span>;
        </div>
    </div>
    <div style="background-color: #1a1210; color: #fef3e6; padding: 20px; border: 2px solid #3d2e26; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #fb923c; margin-top: 0;">Sunset Theme (Dunkel)</h3>
        <p style="color: #b8a090;">Warme Orangen- und Erdtöne, die die goldene Stunde evozieren.</p>
        <div style="background-color: #2d1f1a; border: 1px solid #3d2e26; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #fb923c;">const</span> <span style="color: #fef3e6;">theme</span> = <span style="color: #b8a090;">"sunset"</span>;
        </div>
    </div>
</div>

## Midnight

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #f8f9fc; color: #1e1e2e; padding: 20px; border: 2px solid #ccd0e0; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #5b5fc7; margin-top: 0;">Midnight Theme (Hell)</h3>
        <p style="color: #6c6c8a;">Tiefe Violett- und Blautöne für ein anspruchsvolles, modernes Aussehen.</p>
        <div style="background-color: #eef0f8; border: 1px solid #ccd0e0; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #5b5fc7;">const</span> <span style="color: #1e1e2e;">theme</span> = <span style="color: #6c6c8a;">"midnight"</span>;
        </div>
    </div>
    <div style="background-color: #11111b; color: #cdd6f4; padding: 20px; border: 2px solid #313244; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #89b4fa; margin-top: 0;">Midnight Theme (Dunkel)</h3>
        <p style="color: #9399b2;">Tiefe Violett- und Blautöne für ein anspruchsvolles, modernes Aussehen.</p>
        <div style="background-color: #1e1e2e; border: 1px solid #313244; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #89b4fa;">const</span> <span style="color: #cdd6f4;">theme</span> = <span style="color: #9399b2;">"midnight"</span>;
        </div>
    </div>
</div>

## Lavender

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #faf8ff; color: #2e1e3e; padding: 20px; border: 2px solid #e0d4ec; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #7c3aed; margin-top: 0;">Lavender Theme (Hell)</h3>
        <p style="color: #7c6a8a;">Sanfte Violetttöne und zarte Farben für eine beruhigende, elegante Ästhetik.</p>
        <div style="background-color: #f5f0fa; border: 1px solid #e0d4ec; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #7c3aed;">const</span> <span style="color: #2e1e3e;">theme</span> = <span style="color: #7c6a8a;">"lavender"</span>;
        </div>
    </div>
    <div style="background-color: #1a1625; color: #f0e8f8; padding: 20px; border: 2px solid #3d2e4a; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #a78bfa; margin-top: 0;">Lavender Theme (Dunkel)</h3>
        <p style="color: #a090b8;">Sanfte Violetttöne und zarte Farben für eine beruhigende, elegante Ästhetik.</p>
        <div style="background-color: #2e1e3e; border: 1px solid #3d2e4a; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #a78bfa;">const</span> <span style="color: #f0e8f8;">theme</span> = <span style="color: #a090b8;">"lavender"</span>;
        </div>
    </div>
</div>

## Rose

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fff5f7; color: #2e1a1e; padding: 20px; border: 2px solid #f0d4dc; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #e11d48; margin-top: 0;">Rose Theme (Hell)</h3>
        <p style="color: #8a6c74;">Romantische Rosa- und Rosatöne für ein warmes, einladendes Gefühl.</p>
        <div style="background-color: #fef0f3; border: 1px solid #f0d4dc; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #e11d48;">const</span> <span style="color: #2e1a1e;">theme</span> = <span style="color: #8a6c74;">"rose"</span>;
        </div>
    </div>
    <div style="background-color: #1a1012; color: #fce8ec; padding: 20px; border: 2px solid #4a2e36; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #fb7185; margin-top: 0;">Rose Theme (Dunkel)</h3>
        <p style="color: #b89098;">Romantische Rosa- und Rosatöne für ein warmes, einladendes Gefühl.</p>
        <div style="background-color: #2e1a1e; border: 1px solid #4a2e36; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #fb7185;">const</span> <span style="color: #fce8ec;">theme</span> = <span style="color: #b89098;">"rose"</span>;
        </div>
    </div>
</div>

## Nord

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #eceff4; color: #2e3440; padding: 20px; border: 2px solid #d8dee9; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #5e81ac; margin-top: 0;">Nord Theme (Hell)</h3>
        <p style="color: #6b7280;">Die beliebte Nord-Farbpalette mit arktischen Blautönen und kühlen Grautönen.</p>
        <div style="background-color: #e5e9f0; border: 1px solid #d8dee9; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #5e81ac;">const</span> <span style="color: #2e3440;">theme</span> = <span style="color: #6b7280;">"nord"</span>;
        </div>
    </div>
    <div style="background-color: #2e3440; color: #eceff4; padding: 20px; border: 2px solid #4c566a; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #88c0d0; margin-top: 0;">Nord Theme (Dunkel)</h3>
        <p style="color: #a3a7b1;">Die beliebte Nord-Farbpalette mit arktischen Blautönen und kühlen Grautönen.</p>
        <div style="background-color: #3b4252; border: 1px solid #4c566a; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #88c0d0;">const</span> <span style="color: #eceff4;">theme</span> = <span style="color: #a3a7b1;">"nord"</span>;
        </div>
    </div>
</div>

## Gruvbox

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fbf1c7; color: #3c3836; padding: 20px; border: 2px solid #d5c4a1; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #d65d0e; margin-top: 0;">Gruvbox Theme (Hell)</h3>
        <p style="color: #7c6f64;">Das Retro-Groove-Farbschema mit warmen, gedämpften Tönen.</p>
        <div style="background-color: #f2e5bc; border: 1px solid #d5c4a1; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #d65d0e;">const</span> <span style="color: #3c3836;">theme</span> = <span style="color: #7c6f64;">"gruvbox"</span>;
        </div>
    </div>
    <div style="background-color: #282828; color: #ebdbb2; padding: 20px; border: 2px solid #3c3836; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #fe8019; margin-top: 0;">Gruvbox Theme (Dunkel)</h3>
        <p style="color: #a89984;">Das Retro-Groove-Farbschema mit warmen, gedämpften Tönen.</p>
        <div style="background-color: #1d2021; border: 1px solid #3c3836; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #fe8019;">const</span> <span style="color: #ebdbb2;">theme</span> = <span style="color: #a89984;">"gruvbox"</span>;
        </div>
    </div>
</div>

## Dracula

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #f8f8f2; color: #282a36; padding: 20px; border: 2px solid #dcdcdc; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #bd93f9; margin-top: 0;">Dracula Theme (Hell)</h3>
        <p style="color: #6272a4;">Das beliebte dunkle Theme mit lebendigen Violett- und Rosatönen.</p>
        <div style="background-color: #eff0eb; border: 1px solid #dcdcdc; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #bd93f9;">const</span> <span style="color: #282a36;">theme</span> = <span style="color: #6272a4;">"dracula"</span>;
        </div>
    </div>
    <div style="background-color: #282a36; color: #f8f8f2; padding: 20px; border: 2px solid #44475a; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #bd93f9; margin-top: 0;">Dracula Theme (Dunkel)</h3>
        <p style="color: #6272a4;">Das beliebte dunkle Theme mit lebendigen Violett- und Rosatönen.</p>
        <div style="background-color: #1e1f29; border: 1px solid #44475a; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #bd93f9;">const</span> <span style="color: #f8f8f2;">theme</span> = <span style="color: #6272a4;">"dracula"</span>;
        </div>
    </div>
</div>

## Solarized

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fdf6e3; color: #073642; padding: 20px; border: 2px solid #eee8d5; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #268bd2; margin-top: 0;">Solarized Theme (Hell)</h3>
        <p style="color: #657b83;">Die sorgfältig gestaltete Farbpalette zur Reduzierung von Augenbelastung.</p>
        <div style="background-color: #eee8d5; border: 1px solid #eee8d5; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #268bd2;">const</span> <span style="color: #073642;">theme</span> = <span style="color: #657b83;">"solarized"</span>;
        </div>
    </div>
    <div style="background-color: #002b36; color: #eee8d5; padding: 20px; border: 2px solid #073642; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #268bd2; margin-top: 0;">Solarized Theme (Dunkel)</h3>
        <p style="color: #93a1a1;">Die sorgfältig gestaltete Farbpalette zur Reduzierung von Augenbelastung.</p>
        <div style="background-color: #073642; border: 1px solid #073642; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #268bd2;">const</span> <span style="color: #eee8d5;">theme</span> = <span style="color: #93a1a1;">"solarized"</span>;
        </div>
    </div>
</div>

## Mono

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #ffffff; color: #111111; padding: 20px; border: 2px solid #e5e7eb; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #000000; margin-top: 0;">Mono Theme (Hell)</h3>
        <p style="color: #6b7280;">Ein minimalistisches Schwarz-Weiß-Theme für maximalen Kontrast.</p>
        <div style="background-color: #f9fafb; border: 1px solid #e5e7eb; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #000000;">const</span> <span style="color: #111111;">theme</span> = <span style="color: #6b7280;">"mono"</span>;
        </div>
    </div>
    <div style="background-color: #0a0a0a; color: #f5f5f5; padding: 20px; border: 2px solid #262626; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #ffffff; margin-top: 0;">Mono Theme (Dunkel)</h3>
        <p style="color: #9ca3af;">Ein minimalistisches Schwarz-Weiß-Theme für maximalen Kontrast.</p>
        <div style="background-color: #171717; border: 1px solid #262626; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #ffffff;">const</span> <span style="color: #f5f5f5;">theme</span> = <span style="color: #9ca3af;">"mono"</span>;
        </div>
    </div>
</div>

## Cyberpunk

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fdfcff; color: #1a1025; padding: 20px; border: 2px solid #e0d7ff; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #ff2ea6; margin-top: 0;">Cyberpunk Theme (Hell)</h3>
        <p style="color: #6b5b95;">Neon-Rosa- und Violetttöne für ein futuristisches, High-Tech-Feeling.</p>
        <div style="background-color: #f4f0ff; border: 1px solid #e0d7ff; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #ff2ea6;">const</span> <span style="color: #1a1025;">theme</span> = <span style="color: #6b5b95;">"cyberpunk"</span>;
        </div>
    </div>
    <div style="background-color: #0b0014; color: #f4e9ff; padding: 20px; border: 2px solid #2a104a; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #ff2ea6; margin-top: 0;">Cyberpunk Theme (Dunkel)</h3>
        <p style="color: #a78bfa;">Neon-Rosa- und Violetttöne für ein futuristisches, High-Tech-Feeling.</p>
        <div style="background-color: #1a0033; border: 1px solid #2a104a; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #ff2ea6;">const</span> <span style="color: #f4e9ff;">theme</span> = <span style="color: #a78bfa;">"cyberpunk"</span>;
        </div>
    </div>
</div>

## Desert

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fff8ed; color: #3b2f2f; padding: 20px; border: 2px solid #e6d3b1; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #c0841d; margin-top: 0;">Desert Theme (Hell)</h3>
        <p style="color: #8b7355;">Sandige Beigetöne und warme Goldtöne, inspiriert von trockenen Landschaften.</p>
        <div style="background-color: #fdf1dc; border: 1px solid #e6d3b1; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #c0841d;">const</span> <span style="color: #3b2f2f;">theme</span> = <span style="color: #8b7355;">"desert"</span>;
        </div>
    </div>
    <div style="background-color: #1f1a14; color: #fdf1dc; padding: 20px; border: 2px solid #3d3326; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #fbbf24; margin-top: 0;">Desert Theme (Dunkel)</h3>
        <p style="color: #c9b28a;">Sandige Beigetöne und warme Goldtöne, inspiriert von trockenen Landschaften.</p>
        <div style="background-color: #2a2219; border: 1px solid #3d3326; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #fbbf24;">const</span> <span style="color: #fdf1dc;">theme</span> = <span style="color: #c9b28a;">"desert"</span>;
        </div>
    </div>
</div>

## Ice

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #f0f9ff; color: #0c1e2c; padding: 20px; border: 2px solid #cfe8f3; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #0ea5e9; margin-top: 0;">Ice Theme (Hell)</h3>
        <p style="color: #5b7c99;">Kühle, klare Blautöne, die Winter und Frost evozieren.</p>
        <div style="background-color: #e6f6ff; border: 1px solid #cfe8f3; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #0ea5e9;">const</span> <span style="color: #0c1e2c;">theme</span> = <span style="color: #5b7c99;">"ice"</span>;
        </div>
    </div>
    <div style="background-color: #020617; color: #e0f2fe; padding: 20px; border: 2px solid #082f49; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #38bdf8; margin-top: 0;">Ice Theme (Dunkel)</h3>
        <p style="color: #7dd3fc;">Kühle, klare Blautöne, die Winter und Frost evozieren.</p>
        <div style="background-color: #031525; border: 1px solid #082f49; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #38bdf8;">const</span> <span style="color: #e0f2fe;">theme</span> = <span style="color: #7dd3fc;">"ice"</span>;
        </div>
    </div>
</div>

## Coffee

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #faf6f2; color: #3a2e2a; padding: 20px; border: 2px solid #e4d6cd; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #7c2d12; margin-top: 0;">Coffee Theme (Hell)</h3>
        <p style="color: #7a5e54;">Reiche Brauntöne und warme Farben wie Ihr Lieblingsgetränk.</p>
        <div style="background-color: #f3ebe5; border: 1px solid #e4d6cd; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #7c2d12;">const</span> <span style="color: #3a2e2a;">theme</span> = <span style="color: #7a5e54;">"coffee"</span>;
        </div>
    </div>
    <div style="background-color: #1f1713; color: #f5ede6; padding: 20px; border: 2px solid #3b2a23; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #f97316; margin-top: 0;">Coffee Theme (Dunkel)</h3>
        <p style="color: #b8a29a;">Reiche Brauntöne und warme Farben wie Ihr Lieblingsgetränk.</p>
        <div style="background-color: #2a1e18; border: 1px solid #3b2a23; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #f97316;">const</span> <span style="color: #f5ede6;">theme</span> = <span style="color: #b8a29a;">"coffee"</span>;
        </div>
    </div>
</div>

## Emerald

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #f0fdf4; color: #052e16; padding: 20px; border: 2px solid #bbf7d0; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #059669; margin-top: 0;">Emerald Theme (Hell)</h3>
        <p style="color: #5b8a6e;">Lebendige Grüntöne, die Leben und Energie in Ihre Inhalte bringen.</p>
        <div style="background-color: #dcfce7; border: 1px solid #bbf7d0; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #059669;">const</span> <span style="color: #052e16;">theme</span> = <span style="color: #5b8a6e;">"emerald"</span>;
        </div>
    </div>
    <div style="background-color: #022c22; color: #ecfdf5; padding: 20px; border: 2px solid #064e3b; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #34d399; margin-top: 0;">Emerald Theme (Dunkel)</h3>
        <p style="color: #6ee7b7;">Lebendige Grüntöne, die Leben und Energie in Ihre Inhalte bringen.</p>
        <div style="background-color: #033a2e; border: 1px solid #064e3b; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #34d399;">const</span> <span style="color: #ecfdf5;">theme</span> = <span style="color: #6ee7b7;">"emerald"</span>;
        </div>
    </div>
</div>

## Amber

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fffbeb; color: #3b2f0b; padding: 20px; border: 2px solid #fde68a; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #d97706; margin-top: 0;">Amber Theme (Hell)</h3>
        <p style="color: #8a6d3b;">Goldene Gelbtöne und warme Bernsteintöne für ein helles, fröhliches Aussehen.</p>
        <div style="background-color: #fef3c7; border: 1px solid #fde68a; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #d97706;">const</span> <span style="color: #3b2f0b;">theme</span> = <span style="color: #8a6d3b;">"amber"</span>;
        </div>
    </div>
    <div style="background-color: #1c1402; color: #fef3c7; padding: 20px; border: 2px solid #3f2f05; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #fbbf24; margin-top: 0;">Amber Theme (Dunkel)</h3>
        <p style="color: #facc15;">Goldene Gelbtöne und warme Bernsteintöne für ein helles, fröhliches Aussehen.</p>
        <div style="background-color: #2a1f05; border: 1px solid #3f2f05; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #fbbf24;">const</span> <span style="color: #fef3c7;">theme</span> = <span style="color: #facc15;">"amber"</span>;
        </div>
    </div>
</div>

## Matrix

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #f6fff8; color: #042f1a; padding: 20px; border: 2px solid #c6f6d5; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #16a34a; margin-top: 0;">Matrix Theme (Hell)</h3>
        <p style="color: #4d7c5f;">Digitales Grün auf Dunkel, inspiriert von der ikonischen Sci-Fi-Ästhetik.</p>
        <div style="background-color: #dcfce7; border: 1px solid #c6f6d5; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #16a34a;">const</span> <span style="color: #042f1a;">theme</span> = <span style="color: #4d7c5f;">"matrix"</span>;
        </div>
    </div>
    <div style="background-color: #020f07; color: #d1fae5; padding: 20px; border: 2px solid #134e2a; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #22c55e; margin-top: 0;">Matrix Theme (Dunkel)</h3>
        <p style="color: #4ade80;">Digitales Grün auf Dunkel, inspiriert von der ikonischen Sci-Fi-Ästhetik.</p>
        <div style="background-color: #052e16; border: 1px solid #134e2a; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #22c55e;">const</span> <span style="color: #d1fae5;">theme</span> = <span style="color: #4ade80;">"matrix"</span>;
        </div>
    </div>
</div>

## VS Code Dark

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #ffffff; color: #1e1e1e; padding: 20px; border: 2px solid #e5e7eb; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #007acc; margin-top: 0;">VS Code Dark Theme (Hell)</h3>
        <p style="color: #6b7280;">Vertraute Farben aus Visual Studio Code's dunklem Theme.</p>
        <div style="background-color: #f3f3f3; border: 1px solid #e5e7eb; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #007acc;">const</span> <span style="color: #1e1e1e;">theme</span> = <span style="color: #6b7280;">"vscode-dark"</span>;
        </div>
    </div>
    <div style="background-color: #1e1e1e; color: #d4d4d4; padding: 20px; border: 2px solid #2d2d2d; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #3794ff; margin-top: 0;">VS Code Dark Theme (Dunkel)</h3>
        <p style="color: #9da1a6;">Vertraute Farben aus Visual Studio Code's dunklem Theme.</p>
        <div style="background-color: #252526; border: 1px solid #2d2d2d; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #3794ff;">const</span> <span style="color: #d4d4d4;">theme</span> = <span style="color: #9da1a6;">"vscode-dark"</span>;
        </div>
    </div>
</div>

## Carbon

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #f7f7f7; color: #1a1a1a; padding: 20px; border: 2px solid #dadce0; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #111827; margin-top: 0;">Carbon Theme (Hell)</h3>
        <p style="color: #5f6368;">Tiefe Schwarztöne und Grautöne für ein elegantes, professionelles Erscheinungsbild.</p>
        <div style="background-color: #eeeeee; border: 1px solid #dadce0; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #111827;">const</span> <span style="color: #1a1a1a;">theme</span> = <span style="color: #5f6368;">"carbon"</span>;
        </div>
    </div>
    <div style="background-color: #0f0f0f; color: #e5e5e5; padding: 20px; border: 2px solid #262626; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #f9fafb; margin-top: 0;">Carbon Theme (Dunkel)</h3>
        <p style="color: #9ca3af;">Tiefe Schwarztöne und Grautöne für ein elegantes, professionelles Erscheinungsbild.</p>
        <div style="background-color: #1a1a1a; border: 1px solid #262626; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #f9fafb;">const</span> <span style="color: #e5e5e5;">theme</span> = <span style="color: #9ca3af;">"carbon"</span>;
        </div>
    </div>
</div>

## Sakura

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fff1f2; color: #3f1d2a; padding: 20px; border: 2px solid #fecdd3; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #be185d; margin-top: 0;">Sakura Theme (Hell)</h3>
        <p style="color: #9f5b72;">Zarte Kirschblüten-Rosatöne für eine sanfte, feminine Note.</p>
        <div style="background-color: #ffe4e6; border: 1px solid #fecdd3; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #be185d;">const</span> <span style="color: #3f1d2a;">theme</span> = <span style="color: #9f5b72;">"sakura"</span>;
        </div>
    </div>
    <div style="background-color: #1f0a12; color: #ffe4e6; padding: 20px; border: 2px solid #3f0f1f; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #fb7185; margin-top: 0;">Sakura Theme (Dunkel)</h3>
        <p style="color: #fda4af;">Zarte Kirschblüten-Rosatöne für eine sanfte, feminine Note.</p>
        <div style="background-color: #2a0d18; border: 1px solid #3f0f1f; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #fb7185;">const</span> <span style="color: #ffe4e6;">theme</span> = <span style="color: #fda4af;">"sakura"</span>;
        </div>
    </div>
</div>

## Terminal

<div style="display: flex; gap: 15px; margin: 15px 0; flex-wrap: wrap;">
    <div style="background-color: #fcfcfc; color: #000000; padding: 20px; border: 2px solid #d1d5db; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #16a34a; margin-top: 0;">Terminal Theme (Hell)</h3>
        <p style="color: #6b7280;">Klassische Terminal-Ästhetik mit grünen Akzenten auf Schwarz.</p>
        <div style="background-color: #f4f4f5; border: 1px solid #d1d5db; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #16a34a;">const</span> <span style="color: #000000;">theme</span> = <span style="color: #6b7280;">"terminal"</span>;
        </div>
    </div>
    <div style="background-color: #000000; color: #e5e7eb; padding: 20px; border: 2px solid #1f2937; border-radius: 8px; flex: 1; min-width: 250px;">
        <h3 style="color: #22c55e; margin-top: 0;">Terminal Theme (Dunkel)</h3>
        <p style="color: #9ca3af;">Klassische Terminal-Ästhetik mit grünen Akzenten auf Schwarz.</p>
        <div style="background-color: #020617; border: 1px solid #1f2937; border-radius: 4px; padding: 10px; margin-top: 10px; font-family: monospace; font-size: 0.9em;">
            <span style="color: #22c55e;">const</span> <span style="color: #e5e7eb;">theme</span> = <span style="color: #9ca3af;">"terminal"</span>;
        </div>
    </div>
</div>
