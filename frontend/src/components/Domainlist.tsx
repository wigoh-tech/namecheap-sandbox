// src/components/DomainList.tsx
import { useEffect, useState } from "react";

type Domain = {
  id: number;
  name: string;
 
  
  total?: number;
  revoked: boolean;
  dnsRecord?: {
    aRecord: string;    
    cName: string;
  };
};

export default function DomainList() {
  const [domains, setDomains] = useState<Domain[]>([]);
  const [status, setStatus] = useState("");

  const fetchDomains = async () => {
    const res = await fetch("http://localhost:8080/domains");
    const data = await res.json();

    console.log("Fetched domains:", data);
   setDomains(Array.isArray(data) ? data : data.domains || []);

  };

  const revokeDomain = async (name: string) => {
      console.log("Domain to revoke:", name);


    const res = await fetch("http://localhost:8080/revoke-domain", {
      method: "POST",
      headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ domain: name }), // ✅ Proper body expected by Go backend
  });
    
    const data = await res.json();
    if (res.ok) {
      setStatus("✅ Domain revoked");
      fetchDomains();
    } else {
      setStatus(`❌ ${data.error}`);
    }
  };

  useEffect(() => {
    fetchDomains();
     const listener = () => fetchDomains();
     window.addEventListener("refreshDomains", listener);
     return () => window.removeEventListener("refreshDomains", listener);
  }, []);

  return (
    <div className="mt-8">
      <h2 className="text-2xl font-semibold mb-4">Purchased Domains</h2>
      {status && <p className="text-sm text-gray-700 mb-2">{status}</p>}
      <table className="min-w-full bg-white border shadow rounded-xl">
        <thead>
          <tr className="bg-gray-100 text-left">
            <th className="p-3">Domain</th>
            <th className="p-3">A Record</th>
            <th className="p-3">CNAME</th>
          
            <th className="p-3">Price</th>
            <th className="p-3">Status</th>
            <th className="p-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          {domains.map((d) => (
  <tr key={d.id || d.name}> {/* Use unique key */}
    <td className="p-3">{d.name}</td> {/* Ensure d.name exists */}
    <td className="p-3">{d.dnsRecord?.aRecord ?? "N/A"}</td>
    <td className="p-3">{d.dnsRecord?.cName ?? "N/A"}</td>
    
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
        onClick={() => revokeDomain(d.name)} // ✅ This will now have a valid name
        disabled={d.revoked}
        className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded disabled:opacity-50"
      >
        Revoke
      </button>
      <button
        onClick={() =>
          window.dispatchEvent(new CustomEvent("editDNS", { detail: d }))
        }
        className="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded"
      >
        EDIT DNS
      </button>
    </td>
  </tr>
))}

         
        </tbody>
      </table>
    </div>
  );
}
