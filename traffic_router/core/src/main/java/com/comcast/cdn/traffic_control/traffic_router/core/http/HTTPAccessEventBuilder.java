package com.comcast.cdn.traffic_control.traffic_router.core.http;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.Geolocation;

import javax.servlet.http.HttpServletRequest;
import java.math.RoundingMode;
import java.text.DecimalFormat;
import java.util.Date;
import java.util.Map;

public class HTTPAccessEventBuilder {
    private static String formatRequest(final HttpServletRequest request) {
        String url = formatObject(request.getRequestURL());

        if ("-".equals(url)) {
            return url;
        }

        if (request.getQueryString() != null && !request.getQueryString().isEmpty()) {
            final String queryString = "?" + request.getQueryString();
            final StringBuilder stringBuilder = new StringBuilder(url);
            stringBuilder.append(queryString);
            url = stringBuilder.toString();
        }

        return url;
    }

    private static String formatObject(final Object o) {
        return (o == null) ? "-" : o.toString();
    }

    private static String formatRequestHeaders(final Map<String, String> requestHeaders) {
        if (requestHeaders == null || requestHeaders.isEmpty()) {
            return "rh=\"-\"";
        }

        final StringBuilder stringBuilder = new StringBuilder();
        boolean first = true;
        for (Map.Entry<String, String> entry : requestHeaders.entrySet()) {
            if (entry.getValue() == null || entry.getValue().isEmpty()) {
                continue;
            }

            if (!first) {
                stringBuilder.append(' ');
            }
            else {
                first = false;
            }

            stringBuilder.append("rh=\"");
            stringBuilder.append(entry.getKey()).append(": ");
            stringBuilder.append(entry.getValue().replaceAll("\"", "'"));
            stringBuilder.append('"');
        }

        return stringBuilder.toString();
    }

    @SuppressWarnings("PMD.UseStringBufferForStringAppends")
    public static String create(final HTTPAccessRecord httpAccessRecord) {
        final long start = httpAccessRecord.getRequestDate().getTime();
        final String timeString = String.format("%d.%03d", start / 1000, start % 1000);

        final HttpServletRequest httpServletRequest = httpAccessRecord.getRequest();

        String chi = formatObject(httpServletRequest.getRemoteAddr());
        final String url = formatRequest(httpServletRequest);
        final String cqhm = formatObject(httpServletRequest.getMethod());
        final String cqhv = formatObject(httpServletRequest.getProtocol());

        final String resultType = formatObject(httpAccessRecord.getResultType());
        final String rerr = formatObject(httpAccessRecord.getRerr());

        String resultDetails = "-";
        if (!"-".equals(resultType)) {
            resultDetails = formatObject(httpAccessRecord.getResultDetails());
        }

        String rloc = "-";
        final Geolocation resultLocation = httpAccessRecord.getResultLocation();

        if (resultLocation != null) {
            final DecimalFormat decimalFormat = new DecimalFormat("0.00");
            decimalFormat.setRoundingMode(RoundingMode.DOWN);
            rloc = decimalFormat.format(resultLocation.getLatitude()) + "," + decimalFormat.format(resultLocation.getLongitude());
        }


        final String xMmClientIpHeader = httpServletRequest.getHeader(TRServlet.X_MM_CLIENT_IP);
        final String fakeIpParameter = httpServletRequest.getParameter(TRServlet.FAKE_IP);

        if (xMmClientIpHeader != null) {
            chi = xMmClientIpHeader;
        } else if (fakeIpParameter != null) {
            chi = fakeIpParameter;
        }

        final String rgb = formatObject(httpAccessRecord.getRegionalGeoResult());

        final StringBuilder stringBuilder = new StringBuilder(timeString)
            .append(" qtype=HTTP")
            .append(" chi=" + chi)
            .append(" url=\"" + url + "\"")
            .append(" cqhm=" + cqhm)
            .append(" cqhv=" + cqhv)
            .append(" rtype=" + resultType)
            .append(" rloc=\"" + rloc + "\"")
            .append(" rdtl=" + resultDetails)
            .append(" rerr=\"" + rerr + "\"")
            .append(" rgb=\"" + rgb + "\"");

        if (httpAccessRecord.getResponseCode() != -1) {
            final String pssc = formatObject(httpAccessRecord.getResponseCode());
            final long ttms = new Date().getTime() - start;
            stringBuilder.append(" pssc=").append(pssc).append(" ttms=").append(ttms);
        }

        final String respurl = " rurl=\"" + formatObject(httpAccessRecord.getResponseURL()) + "\" ";
        stringBuilder.append(respurl);

        stringBuilder.append(formatRequestHeaders(httpAccessRecord.getRequestHeaders()));
        return stringBuilder.toString();
    }
}
