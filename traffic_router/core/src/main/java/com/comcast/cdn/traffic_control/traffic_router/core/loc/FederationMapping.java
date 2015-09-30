package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;

import java.util.List;

public class FederationMapping {
    private final String cname;
    private final int ttl;
    private final List<CidrAddress> resolve4;
    private final List<CidrAddress> resolve6;

    public FederationMapping(final String cname, final int ttl, final List<CidrAddress> resolve4, final List<CidrAddress> resolve6) {
        this.cname = cname;
        this.ttl = ttl;
        this.resolve4 = resolve4;
        this.resolve6 = resolve6;
    }

    public String getCname() {
        return cname;
    }

    public int getTtl() {
        return ttl;
    }

    public List<CidrAddress> getResolve4() {
        return resolve4;
    }

    public List<CidrAddress> getResolve6() {
        return resolve6;
    }

    @Override
    @SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity", "PMD.IfStmtsMustUseBraces"})
    public boolean equals(final Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;

        final FederationMapping that = (FederationMapping) o;

        if (ttl != that.ttl) return false;
        if (cname != null ? !cname.equals(that.cname) : that.cname != null) return false;
        if (resolve4 != null ? !resolve4.equals(that.resolve4) : that.resolve4 != null) return false;
        return !(resolve6 != null ? !resolve6.equals(that.resolve6) : that.resolve6 != null);

    }

    @Override
    public int hashCode() {
        int result = cname != null ? cname.hashCode() : 0;
        result = 31 * result + ttl;
        result = 31 * result + (resolve4 != null ? resolve4.hashCode() : 0);
        result = 31 * result + (resolve6 != null ? resolve6.hashCode() : 0);
        return result;
    }

    @Override
    public String toString() {
        return "FederationMapping{" +
                "cname='" + cname + '\'' +
                ", ttl=" + ttl +
                ", resolve4=" + resolve4 +
                ", resolve6=" + resolve6 +
                '}';
    }
}
