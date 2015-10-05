package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import org.xbill.DNS.*;

public class DNSAccessEventBuilder {

    @SuppressWarnings("PMD.UseStringBufferForStringAppends")
    public static String create(final DNSAccessRecord dnsAccessRecord) {
        final String event = createEvent(dnsAccessRecord);
        String rType = "-";
        String rDetails = "-";

        if (dnsAccessRecord.getResultType() != null) {
            rType = dnsAccessRecord.getResultType().toString();
            if (dnsAccessRecord.getResultDetails() != null) {
                rDetails = dnsAccessRecord.getResultDetails().toString();
            }
        }


        final String routingInfo = "rtype=" + rType + " rdetails=" + rDetails + " rerr=\"-\"";
        String answer = "ans=\"-\"";

        if (dnsAccessRecord.getDnsMessage() != null) {
            answer = createAnswer(dnsAccessRecord.getDnsMessage());
        }
        return event + " " + routingInfo + " " + answer;
    }

    private static String createEvent(final DNSAccessRecord dnsAccessRecord) {
        final long finishEpochMillis = System.currentTimeMillis();
        final String timeString = String.format("%d.%03d", dnsAccessRecord.getQueryInstant() / 1000, dnsAccessRecord.getQueryInstant() % 1000);

        final long ttms = finishEpochMillis - dnsAccessRecord.getQueryInstant();
        final String ttmsString = Long.toString(ttms);

        final String addressString = dnsAccessRecord.getClient().getHostAddress();

        final StringBuilder stringBuilder = new StringBuilder(timeString).append(" qtype=DNS").append(" chi=").append(addressString).append(" ttms=").append(ttmsString);

        if (dnsAccessRecord.getDnsMessage() == null) {
            return stringBuilder.append(" xn=- fqdn=- type=- class=- ttl=- rcode=-").toString();
        }

        final String messageHeader = createDnsMessageHeader(dnsAccessRecord.getDnsMessage());
        return stringBuilder.append(messageHeader).toString();
    }

    private static String createDnsMessageHeader(final Message dnsMessage) {
        final String queryHeader = " xn=" + dnsMessage.getHeader().getID();
        final String query = " " + createQuery(dnsMessage.getQuestion());
        final String rcode = " rcode=" + Rcode.string(dnsMessage.getHeader().getRcode());
        return new StringBuilder(queryHeader).append(query).append(rcode).toString();
    }

    private static String createAnswer(final Message dnsMessage) {
        if (dnsMessage.getSectionArray(Section.ANSWER) == null || dnsMessage.getSectionArray(Section.ANSWER).length == 0) {
            return "ans=\"-\"";
        }

        final StringBuilder answerStringBuilder = new StringBuilder();
        for (final Record record : dnsMessage.getSectionArray(Section.ANSWER)) {
            final String s = record.rdataToString() + " ";
            answerStringBuilder.append(s);
        }

        return "ans=\"" + answerStringBuilder.toString().trim() + "\"";
    }

    public static String create(final DNSAccessRecord dnsAccessRecord, final WireParseException wireParseException) {
        final String event = createEvent(dnsAccessRecord);
        final String rerr = "Bad Request:" + wireParseException.getClass().getSimpleName() + ":" + wireParseException.getMessage();
        return new StringBuilder(event)
                .append(" rtype=-")
                .append(" rdetails=-")
                .append(" rerr=\"")
                .append(rerr)
                .append("\"")
                .append(" ans=\"-\"")
                .toString();
    }

    public static String create(final DNSAccessRecord dnsAccessRecord, final Exception exception) {
        final Message dnsMessage = dnsAccessRecord.getDnsMessage();
        dnsMessage.getHeader().setRcode(Rcode.SERVFAIL);
        final String event = createEvent(dnsAccessRecord);

        final String rerr = "Server Error:" + exception.getClass().getSimpleName() + ":" + exception.getMessage();

        return new StringBuilder(event)
                .append(" rtype=-")
                .append(" rdetails=-")
                .append(" rerr=\"")
                .append(rerr)
                .append("\"")
                .append(" ans=\"-\"").toString();
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
