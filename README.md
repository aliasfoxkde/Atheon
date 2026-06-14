<img src="docs/logo.svg" alt="Leakr" width="72" /><br>

# Leakr

![Scanners](https://img.shields.io/badge/scanners-6-blue)
![Java](https://img.shields.io/badge/Java-17%2B-orange)
![Maven Central](https://img.shields.io/maven-central/v/io.github.horadomu/leakr)
![License](https://img.shields.io/badge/license-MIT-green)

**Secret detection for Java applications. A library you embed in your code, and a CLI you run anywhere.**

---

## What is Leakr?

Leakr detects leaked secrets — API keys, tokens, and credentials — in your code, config files, and environment variables. It does one thing and does it well.

It works two ways:

- **As a library**: call `runner.scanString(...)`, `runner.scanFile(...)`, or `runner.scanDir(...)` directly from your Java application. No subprocess. No native binary. Just a dependency.
- **As a CLI tool**: run `leakr scan <path>` from any terminal on any platform with a JRE.

---

## Install

Add to your Maven project:

```xml
<dependency>
    <groupId>io.github.horadomu</groupId>
    <artifactId>leakr</artifactId>
    <version>1.0.0</version>
</dependency>
```

Or with Gradle:

```gradle
implementation 'io.github.horadomu:leakr:1.0.0'
```

---

## Library Usage

```java
import leakr.core.*;
import java.nio.file.*;
import java.util.List;

Runner runner = new Runner(new Registry());

// Scan a string
List<Finding> findings = runner.scanString("AKIAIOSFODNN7EXAMPLE");

// Scan a single file
List<Finding> findings = runner.scanFile(Path.of("config.yaml"));

// Scan a directory (parallel, skips binaries and noise dirs automatically)
List<Finding> findings = runner.scanDir(Path.of("/path/to/project"));

// Scan live environment variables
List<Finding> findings = runner.scanEnv();
```

Each `Finding` exposes: `scanner`, `severity`, `file`, `line`, `description`, `match`.

That's the entire API. No configuration, no setup, no builder pattern.

---

## Why Leakr when gitleaks exists?

[gitleaks](https://github.com/gitleaks/gitleaks) is excellent at what it does: scanning git history for committed secrets. That is not what Leakr does.

|                              | Leakr          | gitleaks       |
|------------------------------|:--------------:|:--------------:|
| **Embeddable Java library**  | **✓**          | **✗**          |
| Scan environment variables   | ✓              |                |
| Scan stdin / arbitrary text  | ✓              | partial        |
| Zero native dependencies     | ✓ (JRE only)   | ✗ (Go binary)  |
| Add a scanner = one class    | ✓              | config + rules |
| Scan files and directories   | ✓              | ✓              |
| Scan git history             |                | ✓              |

If you want git history scanning, use gitleaks. If you want to embed secret detection in a Java application, scan environment variables at startup, or integrate into a pipeline without a binary dependency, use Leakr.

---

## Why Java?

- **Embeddable**: drop it into any Maven or Gradle project and call it from application code, not a shell.
- **One JAR, any platform**: runs on Windows, Linux, and macOS without architecture-specific binaries. If you have a JRE, you have Leakr.
- **Enterprise-native**: Java lives in CI pipelines, build servers, and backend services. Leakr lives there with it.
- **Parallel by default**: directory scans use a thread pool sized to the host CPU automatically.

---

## Built-in Scanners

6 scanners ship out of the box. All fully tested.

| Scanner | Detects | Severity |
|---|---|---|
| `aws-access-key` | AWS access key IDs (AKIA/ASIA…) | Critical |
| `github-pat` | GitHub personal access tokens (ghp_…) | Critical |
| `openai-api-key` | OpenAI API keys (sk-…) | Critical |
| `stripe-secret-key` | Stripe secret keys (sk_live_…) | Critical |
| `slack-bot-token` | Slack bot tokens (xoxb-…) | High |
| `twilio-account-sid` | Twilio account SIDs (AC…) | High |

---

## CLI Quickstart

### Build from source

```bash
git clone https://github.com/HoraDomu/Leakr.git
cd Leakr
mvn package -q
```

This produces two JARs in `target/`:
- `leakr-1.0.0-cli.jar` — self-contained fat JAR for CLI use
- `leakr-1.0.0.jar` — thin library JAR for embedding as a dependency

### Run the CLI

```bash
# Scan a directory
java -jar target/leakr-1.0.0-cli.jar scan /path/to/project

# Scan a single file
java -jar target/leakr-1.0.0-cli.jar scan config.env

# Scan environment variables
java -jar target/leakr-1.0.0-cli.jar scan --env

# Pipe content from another command
git diff | java -jar target/leakr-1.0.0-cli.jar scan --stdin

# JSON output (for CI parsing or downstream tooling)
java -jar target/leakr-1.0.0-cli.jar scan /path/to/project --json

# Exclude directories
java -jar target/leakr-1.0.0-cli.jar scan . --exclude target,dist,node_modules

# Limit to specific file extensions
java -jar target/leakr-1.0.0-cli.jar scan . --ext .env,.yaml,.json,.tf

# List all registered scanners
java -jar target/leakr-1.0.0-cli.jar list
```

Exit code `0` means clean. Exit code `1` means findings were detected.

### Sample output

```
[CRITICAL] aws-access-key
  file:  src/config/dev.properties:14
  desc:  Detects AWS access key IDs (AKIA...)
  match: AKIA**********************MPLE

─────────────────────────────
found 1 potential secret(s)

files: 42  size: 318.7 KB  time: 84ms
```

---

## Adding a Scanner

Leakr auto-discovers every class in the `leakr.scanners` package that implements `Scanner`. No registration, no config file. Drop a class in, rebuild, done.

### 1. Create the class

```java
package leakr.scanners;

import leakr.core.*;
import java.util.*;
import java.util.regex.*;

public class MyServiceScanner implements Scanner {
    private static final Pattern PATTERN = Pattern.compile("myservice_[a-zA-Z0-9]{32}");

    public String name()        { return "myservice-api-key"; }
    public String description() { return "Detects MyService API keys"; }
    public Severity severity()  { return Severity.HIGH; }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);
        while (m.find()) matches.add(m.group());
        return matches;
    }
}
```

Place it in `src/leakr/scanners/MyServiceScanner.java`.

### 2. Rebuild and confirm

```bash
mvn package -q
java -jar target/leakr-1.0.0-cli.jar list
```

Your scanner appears automatically.

### 3. Add a test case

Open `src/leakr/test/ScannerTest.java` and add one entry to `CASES`:

```java
new Case("myservice-api-key",
    "myservice_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",  // must produce a match
    "myservice_tooshort")                            // must NOT produce a match
```

---

## Testing

```bash
java -cp target/leakr-1.0.0-cli.jar leakr.test.ScannerTest
```

```
PASS     aws-access-key
PASS     github-pat
PASS     openai-api-key
PASS     stripe-secret-key
PASS     slack-bot-token
PASS     twilio-account-sid

6 passed, 0 failed
```

Exit code `0` on full pass, `1` on any failure.

---

> [!WARNING]
> Leakr is feature complete. Future releases will be security patches and new scanners only. The library API, CLI, output formats, and exit codes are stable and will not change.

---

## License

MIT — free to use, modify, and distribute. Copyright notice must be preserved. See [LICENSE](LICENSE).
