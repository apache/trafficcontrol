package com.comcast.cdn.traffic_control.traffic_router.core.http;

import javax.servlet.http.HttpServletRequest;
import java.util.Date;

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

    @SuppressWarnings("PMD.UseStringBufferForStringAppends")
    public static String create(final HTTPAccessRecord httpAccessRecord) {
        final long start = httpAccessRecord.getRequestDate().getTime();
        final String timeString = String.format("%d.%03d", start / 1000, start % 1000);

        final HttpServletRequest httpServletRequest = httpAccessRecord.getRequest();

        final String chi = formatObject(httpServletRequest.getRemoteAddr());
        final String url = formatRequest(httpServletRequest);
        final String cqhm = formatObject(httpServletRequest.getMethod());
        final String cqhv = formatObject(httpServletRequest.getProtocol());

        final String resultType = formatObject(httpAccessRecord.getResultType());
        final String rerr = formatObject(httpAccessRecord.getRerr());

        String resultDetails = "-";
        if (!"-".equals(resultType)) {
            resultDetails = formatObject(httpAccessRecord.getResultDetails());
        }

        final StringBuilder stringBuilder = new StringBuilder(timeString)
            .append(" qtype=HTTP")
            .append(" chi=" + chi)
            .append(" url=\"" + url + "\"")
            .append(" cqhm=" + cqhm)
            .append(" cqhv=" + cqhv)
            .append(" rtype=" + resultType)
            .append(" rdetails=" + resultDetails)
            .append(" rerr=\"" + rerr + "\"");

        if (httpAccessRecord.getResponseCode() != -1) {
            final String pssc = formatObject(httpAccessRecord.getResponseCode());
            final long ttms = new Date().getTime() - start;
            stringBuilder.append(" pssc=").append(pssc).append(" ttms=").append(ttms);
        }

        final String respurl = " rurl=\"" + formatObject(httpAccessRecord.getResponseURL()) + "\"";
        stringBuilder.append(respurl);

        return stringBuilder.toString();
    }
}
