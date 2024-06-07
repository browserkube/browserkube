package io.browserkube.selenium.bidirectional_w3c;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.openqa.selenium.By;
import org.openqa.selenium.bidi.LogInspector;
import org.openqa.selenium.bidi.log.BaseLogEntry;
import org.openqa.selenium.bidi.log.JavascriptLogEntry;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

class LogTest extends BaseTest {

    @Test
    void testListenToJavascriptLog()
            throws ExecutionException, InterruptedException, TimeoutException {
        try (LogInspector logInspector = new LogInspector(driver)) {
            CompletableFuture<JavascriptLogEntry> future = new CompletableFuture<>();
            logInspector.onJavaScriptLog(future::complete);

            driver.get("https://www.selenium.dev/selenium/web/bidi/logEntryAdded.html");
            driver.findElement(By.id("jsException")).click();

            JavascriptLogEntry logEntry = future.get(5, TimeUnit.SECONDS);

            Assertions.assertEquals("Error: Not working", logEntry.getText());
            Assertions.assertEquals("javascript", logEntry.getType());
            Assertions.assertEquals(BaseLogEntry.LogLevel.ERROR, logEntry.getLevel());
        }
    }

    @Test
    void testListenToJavascriptErrorLog()
            throws ExecutionException, InterruptedException, TimeoutException {
        try (LogInspector logInspector = new LogInspector(driver)) {
            CompletableFuture<JavascriptLogEntry> future = new CompletableFuture<>();
            logInspector.onJavaScriptException(future::complete);

            driver.get("https://www.selenium.dev/selenium/web/bidi/logEntryAdded.html");
            driver.findElement(By.id("jsException")).click();

            JavascriptLogEntry logEntry = future.get(5, TimeUnit.SECONDS);

            Assertions.assertEquals("Error: Not working", logEntry.getText());
            Assertions.assertEquals("javascript", logEntry.getType());
        }
    }

    @Test
    void testListenToLogsWithMultipleConsumers()
            throws ExecutionException, InterruptedException, TimeoutException {
        try (LogInspector logInspector = new LogInspector(driver)) {
            CompletableFuture<JavascriptLogEntry> completableFuture1 = new CompletableFuture<>();
            logInspector.onJavaScriptLog(completableFuture1::complete);

            CompletableFuture<JavascriptLogEntry> completableFuture2 = new CompletableFuture<>();
            logInspector.onJavaScriptLog(completableFuture2::complete);

            driver.get("https://www.selenium.dev/selenium/web/bidi/logEntryAdded.html");
            driver.findElement(By.id("jsException")).click();

            JavascriptLogEntry logEntry = completableFuture1.get(5, TimeUnit.SECONDS);

            Assertions.assertEquals("Error: Not working", logEntry.getText());
            Assertions.assertEquals("javascript", logEntry.getType());

            logEntry = completableFuture2.get(5, TimeUnit.SECONDS);

            Assertions.assertEquals("Error: Not working", logEntry.getText());
            Assertions.assertEquals("javascript", logEntry.getType());
        }
    }
}