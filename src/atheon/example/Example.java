package atheon.example;

import atheon.core.Finding;
import atheon.core.Registry;
import atheon.core.Runner;

import java.nio.file.Paths;
import java.util.List;

public class Example {

    public static void main(String[] args) throws Exception {
        Registry registry = new Registry();
        Runner runner = new Runner(registry);

        // scan a string
        List<Finding> findings = runner.scanString(
            "export AWS_KEY=AKIAIOSFODNN7EXAMPLE12345678"
        );
        for (Finding f : findings) {
            System.out.printf("[%s] %s — %s%n", f.severity.toUpperCase(), f.scanner, f.match);
        }

        // scan a directory
        List<Finding> dirFindings = runner.scanDir(Paths.get("."));
        System.out.printf("%nscanned directory: %d finding(s)%n", dirFindings.size());
    }
}
