---
title: "Watch-Modus"
description: "Wie Sie den Watch-Modus verwenden"
tags: [watch-mode, development]
date: 2025-12-14
draft: false
---

# Watch-Modus

Sie können Dorcs mit Watch-Modus für die Entwicklung verwenden. Das bedeutet, dass Sie Ihre Änderungen sofort im Browser sehen können, ohne den Server neu starten zu müssen.

Es ist so einfach wie das Hinzufügen des `--watch` Flags zum Befehl beim Starten des Servers.

```bash
./dorcs --watch
```

> [!IMPORTANT]
> Dateiüberwachung über Docker-Volumes funktioniert möglicherweise nicht zuverlässig unter Windows/Mac. Für beste Ergebnisse führen Sie dorcs direkt auf Ihrem Host aus: `./dorcs --dir ./docs --watch`

## Was macht der Watch-Modus?

Der Watch-Modus überwacht das `docs` Verzeichnis auf Änderungen und lädt die Seite automatisch neu. 
Er überwacht auch die `dorcs.yaml` Datei auf Änderungen. Wenn Sie also die Konfiguration ändern, wird der Server automatisch neu geladen.

