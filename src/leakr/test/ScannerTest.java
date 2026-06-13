package leakr.test;

import leakr.core.Registry;
import leakr.core.Scanner;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Universal plug-and-play scanner test.
 *
 * Each Case has three fields:
 *   scanner — the exact string returned by scanner.name()
 *   hit     — a string that MUST produce at least one match
 *   miss    — a string that MUST NOT produce any match
 *
 * To test a new scanner: add one Case to the CASES list, rebuild, run.
 * To remove a scanner: delete its Case.
 * All cases run together — adding yours confirms it works and nothing else regressed.
 *
 * Build:  mvn package -q
 * Run:    java -cp target/leakr.jar leakr.test.ScannerTest
 */
public class ScannerTest {

    record Case(String scanner, String hit, String miss) {}

    private static final List<Case> CASES = List.of(
        new Case("aws-access-key",
            "AKIAIOSFODNN7EXAMPLE",
            "AKIA_TOO_SHORT"),

        new Case("github-pat",
            "ghp_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
            "ghp_tooshort"),

        new Case("stripe-secret-key",
            "sk_live_aaaaaaaaaaaaaaaaaaaaaaaa",
            "sk_live_short"),

        new Case("slack-bot-token",
            "xoxb-12345678901-12345678901-aaaaaaaaaaaaaaaaaaaaaaaa",
            "xoxb-bad-token"),

        new Case("twilio-account-sid",
            "ACaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
            "AC12345"),

        new Case("OpenAI-Api-Key",
                "sk-4u8d9f8g7h6j5k4l3m2n1o0p",
                "sk-313123"
                )
    );

    public static void main(String[] args) {
        Registry registry = new Registry();
        Map<String, Scanner> byName = new HashMap<>();
        for (Scanner s : registry.all()) byName.put(s.name(), s);

        int passed = 0;
        int failed = 0;

        for (Case c : CASES) {
            Scanner s = byName.get(c.scanner);
            if (s == null) {
                System.out.printf("MISSING  %-30s  scanner not registered%n", c.scanner);
                failed++;
                continue;
            }

            boolean hitOk  = !s.scan(c.hit).isEmpty();
            boolean missOk =  s.scan(c.miss).isEmpty();

            if (hitOk && missOk) {
                System.out.printf("PASS     %s%n", c.scanner);
                passed++;
            } else {
                if (!hitOk) System.out.printf("FAIL     %-30s  should match:     %s%n", c.scanner, c.hit);
                if (!missOk) System.out.printf("FAIL     %-30s  should NOT match: %s%n", c.scanner, c.miss);
                failed++;
            }
        }

        System.out.printf("%n%d passed, %d failed%n", passed, failed);
        System.exit(failed > 0 ? 1 : 0);
    }
}
