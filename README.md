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

- **As a library**: call `runner.scanString(...)`, `runner.scanFile(...)`, or `runner.scanDir(...)` directly from your Java application. No subprocess. No native binary. Just a dependency.
- **As a CLI tool**: run `leakr scan <path>` from any terminal on any platform with a JRE.

---

## Why Leakr?

Secrets end up in source code all the time. A developer hardcodes an API key to test something locally, forgets to remove it, and commits it. A `.env` file gets checked in by accident. A config template ships with a real credential still in it. Once a secret is in git history, it is effectively public — even if you delete the file, the commit remains.

Most existing secret scanners are built around git hooks or CI pipelines and are tightly coupled to those workflows. They also tend to be external native tools (`trufflehog`, `gitleaks`, etc.) that you invoke as a subprocess — meaning you cannot embed them inside a Java application, and you cannot scan strings or in-memory content programmatically.

Leakr is designed differently:

**It is a Java library first.** You add it as a Maven or Gradle dependency and call it directly from your code. This makes it useful in scenarios no CLI tool can cover: scanning user-uploaded content before storing it, validating config values at application startup, running secret checks inside a CI step written in Java, or integrating into an internal security tool.

**It ships a CLI too.** For the common case — auditing a directory, scanning a diff, checking env vars — you do not need to write any code. Build once, run `java -jar` anywhere a JRE exists.

**It is zero-config.** There is no YAML config, no rule file, no plugin directory to manage. Drop it in, call it, get results. Scanners are auto-discovered at runtime by package scanning — adding a new one is as simple as writing a class that implements `Scanner`.

**It is honest about its scope.** Leakr uses high-confidence regex patterns targeting well-known secret formats (AWS access keys, GitHub PATs, OpenAI keys, Stripe secret keys, Slack bot tokens, Twilio account SIDs). It does not attempt heuristic or entropy-based detection. That keeps false positives low and results trustworthy.


---

## Install

**Maven**

```xml
<dependency>
    <groupId>io.github.horadomu</groupId>
    <artifactId>leakr</artifactId>
    <version>1.0.1</version>
</dependency>
```

**Gradle**

```gradle
implementation 'io.github.horadomu:leakr:1.0.1'
```

---

## Library Usage

```java
import leakr.core.*;
import java.nio.file.*;
import java.util.List;

Runner runner = new Runner(new Registry());

List<Finding> findings = runner.scanString("AKIAIOSFODNN7EXAMPLE");
List<Finding> findings = runner.scanFile(Path.of("config.yaml"));
List<Finding> findings = runner.scanDir(Path.of("/path/to/project"));
List<Finding> findings = runner.scanEnv();
```

Each `Finding` exposes: `scanner`, `severity`, `file`, `line`, `description`, `match`. No configuration required.

---

## CLI Quickstart

Clone the repo, then build with whichever tool you have available.

**Maven**

```bash
git clone https://github.com/HoraDomu/Leakr.git
cd Leakr
mvn package -q
java -jar target/leakr-1.0.1-cli.jar scan /path/to/project
```

**Gradle**

```bash
git clone https://github.com/HoraDomu/Leakr.git
cd Leakr
./gradlew shadowJar
java -jar build/libs/leakr-1.0.1-cli.jar scan /path/to/project
```

**Plain Java (no build tool)**

```bash
git clone https://github.com/HoraDomu/Leakr.git
cd Leakr
# Pull dependencies into lib/
mvn dependency:copy-dependencies -DoutputDirectory=lib -q

# Compile
javac -cp "lib/*" -d out src/leakr/cli/Main.java src/leakr/core/*.java src/leakr/output/*.java src/leakr/scanners/*.java

# Run
java -cp "out:lib/*" leakr.cli.Main scan /path/to/project
# Windows: use semicolons
java -cp "out;lib/*" leakr.cli.Main scan /path/to/project
```

The Maven `package` goal produces two JARs in `target/`:
- `leakr-1.0.1-cli.jar` — self-contained JAR for CLI use
- `leakr-1.0.1.jar` — thin library JAR for embedding as a dependency

---

## Run the CLI

```bash
# Scan a directory
java -jar target/leakr-1.0.1-cli.jar scan /path/to/project

# Scan a single file
java -jar target/leakr-1.0.1-cli.jar scan config.env

# Scan environment variables
java -jar target/leakr-1.0.1-cli.jar scan --env

# Pipe content from stdin
git diff | java -jar target/leakr-1.0.1-cli.jar scan --stdin

# JSON output
java -jar target/leakr-1.0.1-cli.jar scan /path/to/project --json

# Exclude directories
java -jar target/leakr-1.0.1-cli.jar scan . --exclude target,dist,node_modules

# Filter by file extension
java -jar target/leakr-1.0.1-cli.jar scan . --ext .env,.yaml,.json,.tf

# List registered scanners
java -jar target/leakr-1.0.1-cli.jar list
```

Exit code `0` means clean. Exit code `1` means findings were detected.

---

## Sample Output

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

Leakr auto-discovers every class in the `leakr.scanners` package that implements `Scanner`. No registration, no config file — drop a class in, rebuild, and it is live.

### 1. Create the scanner

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

### 2. Rebuild and verify

```bash
mvn package -q
java -jar target/leakr-1.0.1-cli.jar list
```

### 3. Add a test case

Open `src/leakr/test/ScannerTest.java` and add one entry to `CASES`:

```java
new Case("myservice-api-key",
    "myservice_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",  // must match
    "myservice_tooshort")                            // must not match
```

---

## Testing

```bash
java -cp target/leakr-1.0.1-cli.jar leakr.test.ScannerTest
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

MIT — free to use, modify, and distribute. The copyright notice must be preserved in all copies. See [LICENSE](LICENSE).
