# ContainerShip

Ein einfaches CLI-Tool zur einfachen Verwaltung von Docker-Containern. ContainerShip bietet eine unkomplizierte Schnittstelle zum Versenden, Stoppen, Überwachen und Interagieren mit Ihren containerisierten Anwendungen.

## Funktionen

- **Container versenden**: Stellen Sie Ihre Anwendungen schnell in Containern bereit
- **Container-Verwaltung**: Starten, stoppen und überwachen Sie den Container-Status
- **Logs & Debugging**: Zeigen Sie Container-Logs an und führen Sie Befehle innerhalb von Containern aus
- **Einfache CLI**: Benutzerfreundliche Kommandozeilen-Schnittstelle
- **Docker-Integration**: Basierend auf der Docker-API

## Installation

### Voraussetzungen

- Go 1.24 oder neuer
- Docker installiert und läuft

### Aus dem Quellcode bauen

```bash
git clone https://github.com/Femn0X/ContainerShip.git
cd ContainerShip
make build
```

Die Binärdatei ist unter `bin/cs` verfügbar.

### Global installieren

```bash
make install
```

Dies installiert die Binärdatei in Ihrem GOPATH/bin.

## Verwendung

### Grundlegende Befehle

```bash
# Einen Container versenden (Anwendung bereitstellen)
cs ship

# Laufende Container stoppen
cs stop

# Status der Container überprüfen
cs status

# Definierte Container auflisten
cs list

# Container-Logs anzeigen
cs logs

# Befehl im Container ausführen
cs exec <container_name> <command>
```

### Hilfe

```bash
cs help
```

## Konfiguration

ContainerShip verwendet `containership.yaml` für die Konfiguration. Siehe die Beispieldatei für Details.

## Entwicklung

### Bauen

```bash
make build
```

### Testen

```bash
make test
```

### Aufräumen

```bash
make clean
```

### Ausführen

```bash
make run
```

## Mitwirken

Beiträge sind willkommen! Bitte reichen Sie gerne einen Pull Request ein.

## Lizenz

Dieses Projekt ist unter der MIT-Lizenz lizenziert - siehe die [LICENSE](LICENSE) Datei für Details.

## Änderungsprotokoll

Siehe [CHANGELOG.md](CHANGELOG.md) für die Versionshistorie.
