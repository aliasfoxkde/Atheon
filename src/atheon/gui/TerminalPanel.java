package atheon.gui;

import atheon.core.*;

import javax.swing.*;
import javax.swing.text.*;
import java.awt.*;
import java.nio.file.Path;
import java.util.List;
import java.util.concurrent.ExecutionException;

public class TerminalPanel extends JPanel {

    private static final String VERSION = "v1.0.1";


    private final MainWindow parent;
    private final JTextPane  output;
    private final JTextField input;
    private       Style      styleGreen, styleDim, stylePrompt;
    private       Style      styleCritical, styleHigh, styleMedium, styleLow;

    public TerminalPanel(MainWindow parent) {
        this.parent = parent;
        setLayout(new BorderLayout());
        setBackground(AtheonColors.BG);

        output = buildOutput();
        input  = buildInput();

        JScrollPane scroll = new JScrollPane(output);
        scroll.setBorder(BorderFactory.createEmptyBorder());
        scroll.setBackground(AtheonColors.BG);
        scroll.getVerticalScrollBar().setBackground(new Color(20, 20, 22));
        scroll.getVerticalScrollBar().setForeground(AtheonColors.DIM);

        JPanel inputRow = new JPanel(new BorderLayout());
        inputRow.setBackground(AtheonColors.BG);
        inputRow.setBorder(BorderFactory.createEmptyBorder(2, 10, 8, 10));

        JLabel promptLabel = new JLabel("atheon> ");
        promptLabel.setForeground(AtheonColors.PROMPT);
        promptLabel.setFont(new Font(Font.MONOSPACED, Font.BOLD, 13));

        inputRow.add(promptLabel, BorderLayout.WEST);
        inputRow.add(input,       BorderLayout.CENTER);

        add(scroll,   BorderLayout.CENTER);
        add(inputRow, BorderLayout.SOUTH);

        printHeader();
    }

    // ── setup ────────────────────────────────────────────────────────────────

    private JTextPane buildOutput() {
        JTextPane pane = new JTextPane();
        pane.setEditable(false);
        pane.setBackground(AtheonColors.BG);
        pane.setCaretColor(AtheonColors.GREEN);
        pane.setFont(new Font(Font.MONOSPACED, Font.PLAIN, 13));
        pane.setBorder(BorderFactory.createEmptyBorder(10, 10, 4, 10));

        styleGreen    = style(pane, "green",    AtheonColors.GREEN,    13, false);
        styleDim      = style(pane, "dim",      AtheonColors.DIM,      13, false);
        stylePrompt   = style(pane, "prompt",   AtheonColors.PROMPT,   13, false);
        styleCritical = style(pane, "critical", AtheonColors.CRITICAL, 13, false);
        styleHigh     = style(pane, "high",     AtheonColors.HIGH,     13, false);
        styleMedium   = style(pane, "medium",   AtheonColors.MEDIUM,   13, false);
        styleLow      = style(pane, "low",      AtheonColors.LOW,      13, false);

        return pane;
    }

    private JTextField buildInput() {
        JTextField field = new JTextField();
        field.setBackground(AtheonColors.BG);
        field.setForeground(AtheonColors.GREEN);
        field.setCaretColor(AtheonColors.PROMPT);
        field.setFont(new Font(Font.MONOSPACED, Font.PLAIN, 13));
        field.setBorder(BorderFactory.createEmptyBorder());
        field.addActionListener(e -> {
            String cmd = field.getText().trim();
            field.setText("");
            if (!cmd.isEmpty()) handleCommand(cmd);
        });
        return field;
    }

    private Style style(JTextPane pane, String name, Color color, int size, boolean bold) {
        Style s = pane.addStyle(name, null);
        StyleConstants.setForeground(s, color);
        StyleConstants.setFontFamily(s, Font.MONOSPACED);
        StyleConstants.setFontSize(s, size);
        StyleConstants.setBold(s, bold);
        return s;
    }

    public void focusInput() { input.requestFocusInWindow(); }

    // ── header ───────────────────────────────────────────────────────────────

    private void printHeader() {
        ln("atheon " + VERSION + "  ·  type 'help' for commands", styleDim);
        ln("", styleDim);
    }

    // ── command dispatch ─────────────────────────────────────────────────────

    private void handleCommand(String raw) {
        ln("atheon> " + raw, stylePrompt);
        ln("", styleGreen);

        switch (raw.trim().toLowerCase()) {
            case "scan"      -> cmdScan();
            case "scan file" -> cmdScanFile();
            case "scan env"  -> cmdScanEnv();
            case "list"      -> cmdList();
            case "help"      -> cmdHelp();
            case "clear"     -> cmdClear();
            case "//new"     -> { parent.addTab(); return; }
            case "//exit"    -> System.exit(0);
            default          -> ln("  unknown command: '" + raw + "'  —  type 'help'", styleDim);
        }

        ln("", styleGreen);
    }

    // ── commands ─────────────────────────────────────────────────────────────

    private void cmdScan() {
        JFileChooser fc = new JFileChooser();
        fc.setFileSelectionMode(JFileChooser.DIRECTORIES_ONLY);
        fc.setDialogTitle("choose a folder to scan");
        if (fc.showOpenDialog(this) != JFileChooser.APPROVE_OPTION) {
            ln("  cancelled.", styleDim);
            return;
        }
        Path path = fc.getSelectedFile().toPath();
        ln("  scanning  " + path + " ...", styleDim);
        ln("", styleDim);
        runScan(() -> {
            Runner r = new Runner(new Registry());
            List<Finding> f = r.scanDir(path);
            return new ScanResult(f, r.getStats());
        });
    }

    private void cmdScanFile() {
        JFileChooser fc = new JFileChooser();
        fc.setFileSelectionMode(JFileChooser.FILES_ONLY);
        fc.setDialogTitle("choose a file to scan");
        if (fc.showOpenDialog(this) != JFileChooser.APPROVE_OPTION) {
            ln("  cancelled.", styleDim);
            return;
        }
        Path path = fc.getSelectedFile().toPath();
        ln("  scanning  " + path + " ...", styleDim);
        ln("", styleDim);
        runScan(() -> {
            Runner r = new Runner(new Registry());
            List<Finding> f = r.scanFile(path);
            return new ScanResult(f, r.getStats());
        });
    }

    private void cmdScanEnv() {
        ln("  scanning environment variables ...", styleDim);
        ln("", styleDim);
        runScan(() -> {
            List<Finding> f = new Runner(new Registry()).scanEnv();
            return new ScanResult(f, null);
        });
    }

    private void cmdList() {
        List<Scanner> all = new Registry().all();
        ln(String.format("  registered scanners  (%d)", all.size()), styleGreen);
        ln("", styleGreen);
        ln(String.format("  %-28s  %s", "name", "description"), styleDim);
        ln("  " + "─".repeat(68), styleDim);
        for (Scanner s : all) {
            ln(String.format("  %-28s  %s", s.name(), s.description()), styleGreen);
        }
    }

    private void cmdHelp() {
        ln("  commands:", styleGreen);
        ln("  " + "─".repeat(50), styleDim);
        ln("  scan          scan a folder for leaked secrets", styleGreen);
        ln("  scan file     scan a single file", styleGreen);
        ln("  scan env      scan all environment variables", styleGreen);
        ln("  list          list every registered scanner", styleGreen);
        ln("  clear         clear the terminal", styleGreen);
        ln("  help          show this message", styleGreen);
        ln("", styleGreen);
        ln("  //new         open a second scan session tab", styleDim);
        ln("  //exit        close atheon", styleDim);
    }

    private void cmdClear() {
        try {
            output.getStyledDocument().remove(0, output.getStyledDocument().getLength());
        } catch (BadLocationException ignored) {}
        printHeader();
    }

    // ── scan runner ──────────────────────────────────────────────────────────

    @FunctionalInterface
    interface ScanOp { ScanResult run() throws Exception; }

    record ScanResult(List<Finding> findings, Stats stats) {}

    private void runScan(ScanOp op) {
        input.setEnabled(false);

        new SwingWorker<ScanResult, Void>() {
            @Override protected ScanResult doInBackground() throws Exception { return op.run(); }

            @Override protected void done() {
                try {
                    printFindings(get());
                } catch (ExecutionException e) {
                    Throwable cause = e.getCause() != null ? e.getCause() : e;
                    ln("  error: " + cause.getClass().getSimpleName()
                        + (cause.getMessage() != null ? " — " + cause.getMessage() : ""), styleCritical);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    input.setEnabled(true);
                    input.requestFocusInWindow();
                    ln("", styleGreen);
                }
            }
        }.execute();
    }

    private void printFindings(ScanResult result) {
        List<Finding> findings = result.findings();
        Stats stats = result.stats();

        if (findings.isEmpty()) {
            ln("  no secrets detected", styleGreen);
        } else {
            for (Finding f : findings) {
                Style sev = severityStyle(f.severity);
                ln("  [" + f.severity.toUpperCase() + "]  " + f.scanner, sev);
                if (f.file != null) {
                    String loc = f.line != null ? f.file + ":" + f.line : f.file;
                    ln("    file:   " + loc, styleDim);
                }
                ln("    desc:   " + f.description, styleDim);
                ln("    match:  " + redact(f.match), sev);
                ln("", styleGreen);
            }
            ln("  " + "─".repeat(40), styleDim);
            ln("  found " + findings.size() + " potential secret(s)", styleCritical);
        }

        if (stats != null && stats.files > 0) {
            ln("", styleDim);
            ln("  files: " + stats.files
                + "  ·  size: " + formatBytes(stats.bytes)
                + "  ·  time: " + stats.elapsedMs + "ms", styleDim);
        }
    }

    // ── helpers ──────────────────────────────────────────────────────────────

    private Style severityStyle(String sev) {
        return switch (sev) {
            case "critical" -> styleCritical;
            case "high"     -> styleHigh;
            case "medium"   -> styleMedium;
            default         -> styleLow;
        };
    }

    private void ln(String text, Style s) {
        try {
            StyledDocument doc = output.getStyledDocument();
            doc.insertString(doc.getLength(), text + "\n", s);
            output.setCaretPosition(doc.getLength());
        } catch (BadLocationException ignored) {}
    }

    private static String redact(String s) {
        if (s == null || s.length() <= 8) return s == null ? "" : "*".repeat(s.length());
        return s.substring(0, 4) + "*".repeat(s.length() - 8) + s.substring(s.length() - 4);
    }

    private static String formatBytes(long b) {
        if (b >= 1L << 20) return String.format("%.1f MB", (double) b / (1L << 20));
        if (b >= 1L << 10) return String.format("%.1f KB", (double) b / (1L << 10));
        return b + " B";
    }
}
