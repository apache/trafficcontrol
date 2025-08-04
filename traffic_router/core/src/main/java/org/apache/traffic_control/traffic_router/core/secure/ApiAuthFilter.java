package org.apache.traffic_control.traffic_router.core.security;

import javax.servlet.*;
import javax.servlet.http.*;
import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.*;
import java.util.concurrent.atomic.AtomicLong;
import java.util.Base64;

public class ApiAuthFilter implements Filter {

    private File userFile;
    private final Map<String, String> users = new HashMap<>();
    private long lastModified = 0;

    public void setUserFile(final File userFile) {
        this.userFile = userFile;
    }

    @SuppressWarnings("PMD.AvoidThrowingRawExceptionTypes")
    private synchronized void loadUsersIfModified() {
        if (userFile.lastModified() > lastModified) {
            final Map<String, String> tempUsers = new HashMap<>();
            try (BufferedReader reader = new BufferedReader(new FileReader(userFile))) {
                String line;
                while ((line = reader.readLine()) != null) {
                    if (!line.trim().isEmpty() && line.contains(":")) {
                        final String[] parts = line.split(":", 2);
                        tempUsers.put(parts[0], parts[1]);
                    }
                }
                users.clear();
                users.putAll(tempUsers);
                lastModified = userFile.lastModified();
            } catch (IOException e) {
                throw new RuntimeException("Failed to load users from file: " + userFile.getAbsolutePath(), e);
            }
        }
    }

    @Override
    public void doFilter(final ServletRequest req, final ServletResponse res, final FilterChain chain)
        throws IOException, ServletException {

        final HttpServletRequest request = (HttpServletRequest) req;
        final HttpServletResponse response = (HttpServletResponse) res;

        loadUsersIfModified();

        final String authHeader = request.getHeader("Authorization");
        if (authHeader == null || !authHeader.startsWith("Basic ")) {
            response.setHeader("WWW-Authenticate", "Basic realm=\"TrafficRouter\"");
            response.sendError(HttpServletResponse.SC_UNAUTHORIZED);
            return;
        }

        final String base64Credentials = authHeader.substring("Basic ".length());
        final String credentials = new String(Base64.getDecoder().decode(base64Credentials), StandardCharsets.UTF_8);
        final String[] values = credentials.split(":", 2);
        final String username = values[0];
        final String password = values.length > 1 ? values[1] : "";

        if (!users.containsKey(username) || !users.get(username).equals(password)) {
            response.sendError(HttpServletResponse.SC_UNAUTHORIZED);
            return;
        }

        chain.doFilter(req, res);
    }

    @Override public void init(final FilterConfig filterConfig) {}
    @Override public void destroy() {}
}
