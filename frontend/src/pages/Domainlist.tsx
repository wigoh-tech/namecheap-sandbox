// src/pages/DomainList.tsx
import { useEffect, useState } from "react";

type DNSRecord = {
  type: string;
  host: string;
  value: string;
  ttl: number;
};

type Domain = {
  id: number;
  name: string;
  total: number;
  revoked: boolean;
  dnsRecords?: DNSRecord[];
};

export default function DomainList() {
  const [domains, setDomains] = useState<Domain[]>([]);

  const fetchDomains = async () => {
    try {
      const res = await fetch("http://localhost:8080/domains");
      const data = await res.json();
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
    <div className="mt-8 px-4">
      <h2 className="text-3xl font-semibold mb-6 text-gray-800">Purchased Domains</h2>

      <div className="overflow-x-auto bg-white border rounded-lg shadow-md">
        <table className="min-w-full text-sm text-left text-gray-700">
          <thead className="bg-gray-100 text-gray-800 text-base">
            <tr>
              <th className="py-3 px-5">Domain</th>
              <th className="py-3 px-5">Records</th>
              <th className="py-3 px-5">Price</th>
              <th className="py-3 px-5">Status</th>
              <th className="py-3 px-5 text-center">Actions</th>
            </tr>
          </thead>
          <tbody>
            {domains.map((d) => (
              <tr key={d.id} className="border-t hover:bg-gray-50">
                <td className="py-4 px-5 font-medium">{d.name}</td>

                <td className="py-4 px-5 text-gray-600">
                  {Array.isArray(d.dnsRecords) && d.dnsRecords.length > 0 ? (
                    d.dnsRecords.map((r, i) => (
                      <div key={i} className="mb-1">
                        <strong>{r.type}</strong> | {r.host} → {r.value}{" "}
                        <span className="text-xs text-gray-400">(TTL: {r.ttl})</span>
                      </div>
                    ))
                  ) : (
                    <span className="italic text-gray-400">No records</span>
                  )}
                </td>

                <td className="py-4 px-5 font-semibold">₹{d.total?.toFixed(2) ?? "N/A"}</td>

                <td className="py-4 px-5">
                  {d.revoked ? (
                    <span className="text-red-600 font-semibold">Revoked</span>
                  ) : (
                    <span className="text-green-600 font-semibold">Active</span>
                  )}
                </td>

                <td className="py-4 px-5 text-center space-y-2">
                  <button
                    disabled={d.revoked}
                    onClick={() => revokeDomain(d.name)}
                    className="bg-pink-600 hover:bg-yellow-700 text-white px-3 py-1 rounded disabled:opacity-50 w-full"
                  >
                    Revoke
                  </button>

                  <button
                    onClick={() =>
                      window.dispatchEvent(new CustomEvent("editDNS", { detail: { name: d.name } }))
                    }
                    className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded w-full"
                  >
                    Edit DNS
                  </button>

                  <button
                    onClick={() =>
                      window.dispatchEvent(new CustomEvent("addDNS", { detail: { name: d.name } }))
                    }
                    className="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded w-full"
                  >
                    Add DNS
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
