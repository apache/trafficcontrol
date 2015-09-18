package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import org.xbill.DNS.*;

import java.net.InetAddress;

public class DNSAccessEventBuilder {

    private static String create(final long startEpochMillis, final InetAddress clientAddress) {
        final long finishEpochMillis = System.currentTimeMillis();
        final String timeString = String.format("%d.%03d", startEpochMillis / 1000, startEpochMillis % 1000);

        final long ttms = finishEpochMillis - startEpochMillis;
        final String ttmsString = Long.toString(ttms);

        final String addressString = clientAddress.getHostAddress();

        return new StringBuilder(timeString).append(" qtype=DNS").append(" chi=").append(addressString).append(" ttms=").append(ttmsString).toString();
    }

    public static String create(final long startEpochMillis, final InetAddress clientAddress, final WireParseException wireParseException) {
        final String event = create(startEpochMillis, clientAddress);
        final String emptyQuery = " " + createEmptyQuery();
        final String answer = "Bad Request:" + wireParseException.getClass().getSimpleName() + ":" + wireParseException.getMessage();
        return new StringBuilder(event)
                .append(emptyQuery)
                .append(" ans=\"")
                .append(answer)
                .append("\"").toString();
    }

    public static String create(final long startEpochMillis, final InetAddress clientAddress, final Message dnsMessage, final Exception exception) {
        final String event = create(startEpochMillis, clientAddress);
        final String queryHeader = " xn=" + dnsMessage.getHeader().getID();
        final String query = " " + createQuery(dnsMessage.getQuestion());
        final String responseHeader = " rcode=" + Rcode.string(Rcode.SERVFAIL);

        final String answer = "Server Error:" + exception.getClass().getSimpleName() + ":" + exception.getMessage();

        return new StringBuilder(event).append(queryHeader).append(query).append(responseHeader)
                .append(" ans=\"")
                .append(answer)
                .append("\"").toString();
    }

    public static String create(final long startEpochMillis, final InetAddress clientAddress, final Message dnsMessage) {
        final String event = create(startEpochMillis, clientAddress);
        final String queryHeader = " xn=" + dnsMessage.getHeader().getID();
        final String query = " " + createQuery(dnsMessage.getQuestion());
        final String responseHeader = " rcode=" + Rcode.string(dnsMessage.getHeader().getRcode());

        final StringBuilder stringBuilder = new StringBuilder(event).append(queryHeader).append(query).append(responseHeader);


        if (dnsMessage.getSectionArray(Section.ANSWER) == null || dnsMessage.getSectionArray(Section.ANSWER).length == 0) {
            stringBuilder.append(" ans=\"-");
        }
        else {
            stringBuilder.append(" ans=\"");
            for (final Record record : dnsMessage.getSectionArray(Section.ANSWER)) {
                final String s = record.rdataToString() + " ";
                stringBuilder.append(s);
            }
        }

        return stringBuilder.toString().trim() + "\"";
    }

    private static String createEmptyQuery() {
        return "xn=- fqdn=- type=- class=- ttl=- rcode=-";
    }

    private static String createQuery(final Record query) {
        final String qname = query.getName().toString();
        final String qtype = Type.string(query.getType());
        final String qclass = DClass.string(query.getDClass());
        final long ttl = query.getTTL();

        return new StringBuilder()
            .append("fqdn=").append(qname)
            .append(" type=").append(qtype)
            .append(" class=").append(qclass)
            .append(" ttl=").append(ttl)
            .toString();
    }

}
