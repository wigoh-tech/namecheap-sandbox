import { useEffect, useState } from "react";

interface FormState {
  domain: string;
  recordType: string;
  host: string;
  value: string;
  ttl: number;
}

type Mode = "edit" | "add";

export default function DNSModal() {
  const [visible, setVisible] = useState(false);
  const [mode, setMode] = useState<Mode>("edit");
  const [form, setForm] = useState<FormState>({
    domain: "",
    recordType: "A",
    host: "@",
    value: "",
    ttl: 1800,
  });

  useEffect(() => {
    const editHandler = (e: CustomEvent) => {
      const d = e.detail;
      if (!d.name) return alert("Missing domain name");
      setMode("edit");
      setForm({
        domain: d.name,
        recordType: "A",
        host: "@",
        value: "",
        ttl: 1800,
      });
      setVisible(true);
    };

    const addHandler = (e: CustomEvent) => {
      const d = e.detail;
      if (!d.name) return alert("Missing domain name");
      setMode("add");
      setForm({
        domain: d.name,
        recordType: "A",
        host: "@",
        value: "",
        ttl: 1800,
      });
      setVisible(true);
    };

    window.addEventListener("editDNS", editHandler as EventListener);
    window.addEventListener("addDNS", addHandler as EventListener);
    return () => {
      window.removeEventListener("editDNS", editHandler as EventListener);
      window.removeEventListener("addDNS", addHandler as EventListener);
    };
  }, []);

  useEffect(() => {
    const esc = (e: KeyboardEvent) => {
      if (e.key === "Escape") setVisible(false);
    };
    window.addEventListener("keydown", esc);
    return () => window.removeEventListener("keydown", esc);
  }, []);

  const handleSubmit = async () => {
    if (!form.domain || !form.host || !form.value || !form.recordType) {
      alert("All fields are required.");
      return;
    }

    const url =
      mode === "add"
        ? "http://localhost:8080/api/add-dns-record"
        : "http://localhost:8080/api/update-dns";

    try {
      const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(form),
      });

      const text = await res.text();

      if (!res.ok) {
        alert("Failed: " + text);
        return;
      }

      console.log("✅ Success:", text);
      setVisible(false);
      window.dispatchEvent(new Event("refreshDomains"));
    } catch (err) {
      console.error("❌ Error:", err);
      alert("Network or server error.");
    }
  };

  if (!visible) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm">
      <div className="bg-white rounded-lg p-6 shadow-lg w-full max-w-md">
        <h2 className="text-xl font-bold mb-4 text-gray-800">
          {mode === "add" ? "Add" : "Edit"} DNS for{" "}
          <span className="text-blue-600">{form.domain}</span>
        </h2>

        <div className="space-y-3">
          <select
            value={form.recordType}
            onChange={(e) =>
              setForm((f) => ({ ...f, recordType: e.target.value }))
            }
            className="w-full border px-3 py-2 rounded"
          >
            <option>A</option>
            <option>AAAA</option>
            <option>CNAME</option>
            <option>TXT</option>
            <option>DYNAMIC</option>
            <option>ALIAS</option>
            <option>URL</option>
            <option>CAA</option>
          </select>

          <input
            type="text"
            value={form.host}
            onChange={(e) =>
              setForm((f) => ({ ...f, host: e.target.value }))
            }
            placeholder="Host (e.g., @, www)"
            className="w-full border px-3 py-2 rounded"
          />

          <input
            type="text"
            value={form.value}
            onChange={(e) =>
              setForm((f) => ({ ...f, value: e.target.value }))
            }
            placeholder="Value (IP or URL)"
            className="w-full border px-3 py-2 rounded"
          />

          <input
            type="number"
            value={form.ttl}
            onChange={(e) =>
              setForm((f) => ({ ...f, ttl: parseInt(e.target.value) }))
            }
            placeholder="TTL (e.g., 1800)"
            className="w-full border px-3 py-2 rounded"
          />

          <div className="flex justify-end gap-4 pt-4">
            <button
              onClick={() => setVisible(false)}
              className="px-4 py-2 bg-gray-300 hover:bg-gray-400 rounded"
            >
              Cancel
            </button>
            <button
              onClick={handleSubmit}
              className="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded"
            >
              {mode === "add" ? "Add DNS Record" : "Update DNS"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
