import React, { useState } from "react";

interface Props {
  domain: string;
  currentA: string;
  currentCNAME: string;
  onClose: () => void;
  onUpdated: () => void;
}

const EditDNSModal: React.FC<Props> = ({ domain, currentA, currentCNAME, onClose, onUpdated }) => {
  const [aRecord, setARecord] = useState(currentA);
  const [cname, setCname] = useState(currentCNAME);
  const [loading, setLoading] = useState(false);

  const handleUpdate = async () => {
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/move-domain", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ domain, aRecord, cName: cname }),
      });

      const result = await res.json();
      if (res.ok) {
        alert("✅ DNS updated");
        onUpdated(); // refresh list
        onClose();
      } else {
        alert("❌ Update failed: " + result.error);
      }
    } catch (err) {
      alert("❌ Network error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded w-96">
        <h2 className="text-lg font-semibold mb-4">Edit DNS for {domain}</h2>

        <label className="block mb-2 text-sm">A Record</label>
        <input
          value={aRecord}
          onChange={(e) => setARecord(e.target.value)}
          className="w-full p-2 border mb-4"
        />

        <label className="block mb-2 text-sm">CNAME</label>
        <input
          value={cname}
          onChange={(e) => setCname(e.target.value)}
          className="w-full p-2 border mb-4"
        />

        <div className="flex justify-end space-x-2">
          <button onClick={onClose} className="px-4 py-1 bg-gray-300 rounded">Cancel</button>
          <button onClick={handleUpdate} className="px-4 py-1 bg-blue-600 text-white rounded">
            {loading ? "Updating..." : "Update"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default EditDNSModal;
