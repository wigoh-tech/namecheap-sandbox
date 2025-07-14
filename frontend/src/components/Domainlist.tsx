import { useEffect, useState } from "react";

type DNSRecord = {
  type: string;
  host: string;
  value: string;
  ttl: number;
  mxPref?: number;
  flag?: number;
  tag?: string;
};

type Domain = {
  id: number;
  name: string;
  total: number;
  revoked: boolean;
  dnsRecords?: DNSRecord[]; // ✅ optional
};

export default function DomainList() {
  const [domains, setDomains] = useState<Domain[]>([]);

  const fetchDomains = async () => {
    try {
      const res = await fetch("http://localhost:8080/domains");
      const data = await res.json();
      // Ensure data is array
      setDomains(Array.isArray(data) ? data : data.domains || []);
    } catch (err) {
      console.error("Failed to fetch domains:", err);
    }
  };

  const revokeDomain = async (name: string) => {
    await fetch("http://localhost:8080/revoke-domain", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ domain: name }),
    });
    fetchDomains();
  };

  useEffect(() => {
    fetchDomains();
    const listener = () => fetchDomains();
    window.addEventListener("refreshDomains", listener);
    return () => window.removeEventListener("refreshDomains", listener);
  }, []);

  return (
    <div className="mt-8">
      <h2 className="text-2xl font-bold mb-4">Purchased Domains</h2>
      <table className="min-w-full bg-white border shadow rounded-xl">
        <thead>
          <tr className="bg-gray-100 text-left">
            <th className="p-3">Domain</th>
            <th className="p-3">Records</th>
            <th className="p-3">Price</th>
            <th className="p-3">Status</th>
            <th className="p-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          {domains.map((d) => (
            <tr key={d.id}>
              <td className="p-3">{d.name}</td>

              {/* DNS Records column */}
              <td className="p-3 text-sm">
                {Array.isArray(d.dnsRecords) && d.dnsRecords.length > 0 ? (
                  d.dnsRecords.map((r, i) => (
                    <div key={i}>
                      <strong>{r.type}</strong> | {r.host} → {r.value} (TTL: {r.ttl})
                    </div>
                  ))
                ) : (
                  <span className="text-gray-400 italic">No records</span>
                )}
              </td>

              <td className="p-3">₹{d.total?.toFixed(2) ?? "N/A"}</td>

              <td className="p-3">
                {d.revoked ? (
                  <span className="text-red-600 font-semibold">Revoked</span>
                ) : (
                  <span className="text-green-600 font-semibold">Active</span>
                )}
              </td>

              <td className="p-3 space-x-2">
                <button
                  disabled={d.revoked}
                  onClick={() => revokeDomain(d.name)}
                  className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded disabled:opacity-50"
                >
                  Revoke
                </button>

                <button
                  onClick={() =>
                    window.dispatchEvent(
                      new CustomEvent("editDNS", {
                        detail: {
                          domain: d.name,
                          dnsRecords: d.dnsRecords ?? [],
                        },
                      })
                    )
                  }
                  className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded"
                >
                  Edit DNS
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
