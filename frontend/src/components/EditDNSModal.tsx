import React, { useState, useEffect } from "react";

interface Domain {
  domain: string;
  recordType: string;
  aRecord: string;
  cName: string;
}

const EditDNSModal: React.FC = () => {
  const [visible, setVisible] = useState(false);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState<Domain>({
    domain: "",
    recordType: "A", // default to A record
    aRecord: "",
    cName: "",
  });

  const handleClose = () => setVisible(false);

  const handleUpdate = async () => {
    setLoading(true);
    try {
      const body: any = {
        domain: form.domain,
        recordType: form.recordType,
      };

      if (form.recordType === "A") body.aRecord = form.aRecord;
      if (form.recordType === "CNAME") body.cName = form.cName;

      const res = await fetch("http://localhost:8080/api/update-dns", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      const result = await res.json();
      if (res.ok) {
        alert("✅ DNS updated");
        window.dispatchEvent(new CustomEvent("refreshDomains"));
        handleClose();
      } else {
        alert("❌ Update failed: " + result.error);
      }
    } catch (err) {
      alert("❌ Network error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const handler = (e: Event) => {
      const custom = e as CustomEvent<{
        name: string;
        aRecord: string;
        cName: string;
      }>;
      setForm({
        domain: custom.detail.name,
        recordType: "A",
        aRecord: custom.detail.aRecord,
        cName: custom.detail.cName,
      });
      setVisible(true);
    };
    window.addEventListener("editDNS", handler);
    return () => window.removeEventListener("editDNS", handler);
  }, []);

  if (!visible) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded w-96 shadow-xl">
        <h2 className="text-lg font-semibold mb-4">
          Edit DNS for <span className="text-blue-600">{form.domain}</span>
        </h2>

        {/* Record Type Dropdown */}
        <label className="block mb-2 text-sm">Record Type</label>
        <select
          value={form.recordType}
          onChange={(e) => setForm((f) => ({ ...f, recordType: e.target.value }))}
          className="w-full p-2 border mb-4 rounded"
        >
          <option value="A">A Record</option>
          <option value="CNAME">CNAME Record</option>
          {/* You can add more later like MX, NS, etc */}
        </select>

        {/* Conditionally render A Record input */}
        {form.recordType === "A" && (
          <>
            <label className="block mb-2 text-sm">A Record (IP Address)</label>
            <input
              value={form.aRecord}
              onChange={(e) =>
                setForm((f) => ({ ...f, aRecord: e.target.value }))
              }
              className="w-full p-2 border mb-4 rounded"
              placeholder="e.g., 82.25.106.75"
            />
          </>
        )}

        {/* Conditionally render CNAME input */}
        {form.recordType === "CNAME" && (
          <>
            <label className="block mb-2 text-sm">CNAME (Host)</label>
            <input
              value={form.cName}
              onChange={(e) =>
                setForm((f) => ({ ...f, cName: e.target.value }))
              }
              className="w-full p-2 border mb-4 rounded"
              placeholder="e.g., example.hostingersite.com"
            />
          </>
        )}

        <div className="flex justify-end space-x-2">
          <button
            onClick={handleClose}
            className="px-4 py-1 bg-gray-300 hover:bg-gray-400 rounded"
          >
            Cancel
          </button>
          <button
            onClick={handleUpdate}
            className="px-4 py-1 bg-blue-600 text-white rounded"
          >
            {loading ? "Updating..." : "Update"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default EditDNSModal;
