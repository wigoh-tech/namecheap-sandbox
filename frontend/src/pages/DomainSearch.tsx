import { useState } from "react";

type Props = {
  onAvailable: (domain: string) => void;
};

export default function DomainSearch({ onAvailable }: Props) {
  const [domain, setDomain] = useState("");
  const [status, setStatus] = useState("");

  const checkAvailability = async () => {
    try {
      const res = await fetch(`http://localhost:8080/check-domain?domain=${domain}`);

      const data = await res.json();
      if (data.available) {
        setStatus("✅ Domain is available!");
        onAvailable(domain);
      } else {
        setStatus("❌ Domain is already taken.");
      }
    } catch {
      setStatus("❌ Failed to check domain.");
    }
  };

  return (
    <div className="space-y-3">
      <input
        type="text"
        value={domain}
        onChange={(e) => setDomain(e.target.value)}
        placeholder="e.g. mydomain123.com"
        className="w-full px-4 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <button
        onClick={checkAvailability}
        className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition duration-200"
      >
        Check Availability
      </button>
      {status && <p className="text-sm text-gray-700">{status}</p>}
    </div>
  );
}
