package com.fancyinnovations.fancyspaces.utils;

import com.google.gson.Gson;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpResponse;
import java.net.http.HttpTimeoutException;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

public class HttpRequest {

    public static final Gson gson = new Gson();
    private static final Executor executor = Executors.newSingleThreadExecutor();
    private static final HttpClient client = HttpClient.newBuilder()
            .executor(executor)
            .connectTimeout(Duration.ofSeconds(5))
            .build();

    private static AtomicInteger timeoutCount = new AtomicInteger(0);
    private static AtomicLong pauseRequestsUntil = new AtomicLong(0);

    private final String url;
    private String method = "GET";
    private Object body = null;
    private Map<String, String> headers = new HashMap<>();
    private Duration timeout = Duration.ofSeconds(5);
    private HttpResponse.BodyHandler<?> bodyHandler = HttpResponse.BodyHandlers.ofString();

    /**
     * Constructs a new HttpRequest with the provided URL, method, body, and headers.
     *
     * @param url     The URL to which the request is directed.
     * @param method  The HTTP method to use (e.g., GET, POST).
     * @param body    The body of the request. This can be any object that needs to be sent in the request body.
     * @param headers A map representing the headers to be included in the request, where the key is
     *                the header name and the value is the header value.
     */
    public HttpRequest(String url, String method, Object body, Map<String, String> headers) {
        this.url = url;
        this.method = method;
        this.body = body;
        this.headers = headers;
    }

    public HttpRequest(String url) {
        this.url = url;
    }

    /**
     * Sends an HTTP request based on preset configuration including URL, method, body, and headers.
     *
     * @return the HttpResponse containing the response as a string from the executed HTTP request.
     * @throws URISyntaxException   if the URL is not properly formatted.
     * @throws IOException          if an I/O error occurs when sending or receiving.
     * @throws InterruptedException if the operation is interrupted.
     */
    public<T> HttpResponse<T> send() throws URISyntaxException, IOException, InterruptedException {
        if (System.currentTimeMillis() < pauseRequestsUntil.get()) {
            throw new HttpTimeoutException("Request paused due to previous timeouts.");
        }

        URI uri = new URI(url);

        java.net.http.HttpRequest.Builder builder = java.net.http.HttpRequest.newBuilder()
                .uri(uri)
                .timeout(timeout)
                .header("User-Agent", "FancyAnalytics Java-SDK");

        if (!method.equalsIgnoreCase("GET") && body != null) {
            String json = gson.toJson(body);
            if (!json.isEmpty()) {
                builder.method(method, java.net.http.HttpRequest.BodyPublishers.ofString(json));
            }
        } else {
            builder.method(method, java.net.http.HttpRequest.BodyPublishers.noBody());
        }

        if (!headers.isEmpty()) {
            for (Map.Entry<String, String> entry : headers.entrySet()) {
                builder.header(entry.getKey(), entry.getValue());
            }
        }

        try {
            return (HttpResponse<T>) client.send(builder.build(), bodyHandler);
        } catch (HttpTimeoutException e) {
            int timeouts = timeoutCount.incrementAndGet();
            if (timeouts >= 3) {
                long pauseUntil = System.currentTimeMillis() + 1000 * 60 * 60; // Pause for 1 hour
                pauseRequestsUntil.set(pauseUntil);
                timeoutCount.set(0); // Reset timeout count after pausing
            }
            throw e;
        }
    }

    public String getUrl() {
        return url;
    }

    public String getMethod() {
        return method;
    }

    /**
     * Sets the method of the request.
     *
     * @param method The method to use.
     * @return The HttpRequest.
     */
    public HttpRequest withMethod(String method) {
        this.method = method;
        return this;
    }

    public Object getBody() {
        return body;
    }

    /**
     * Sets the body of the request.
     *
     * @param body The body of the request.
     * @return The HttpRequest.
     */
    public HttpRequest withBody(Object body) {
        this.body = body;
        return this;
    }

    public Map<String, String> getHeaders() {
        return headers;
    }

    /**
     * Sets the headers of the request.
     *
     * @param headers The headers of the request.
     * @return The HttpRequest.
     */
    public HttpRequest withHeaders(Map<String, String> headers) {
        this.headers = headers;
        return this;
    }

    /**
     * Adds a header to the request.
     *
     * @param key   The key of the header.
     * @param value The value of the header.
     * @return The HttpRequest.
     */
    public HttpRequest withHeader(String key, String value) {
        this.headers.put(key, value);
        return this;
    }

    public Duration getTimeout() {
        return timeout;
    }

    /**
     * Sets the timeout of the request.
     *
     * @param timeout The timeout of the request.
     * @return The HttpRequest.
     */
    public HttpRequest withTimeout(Duration timeout) {
        this.timeout = timeout;
        return this;
    }

    public HttpResponse.BodyHandler<?> getBodyHandler() {
        return bodyHandler;
    }

    /**
     * Sets the body handler of the request.
     *
     * @param bodyHandler The body handler of the request.
     * @return The HttpRequest.
     */
    public HttpRequest withBodyHandler(HttpResponse.BodyHandler<?> bodyHandler) {
        this.bodyHandler = bodyHandler;
        return this;
    }
}
