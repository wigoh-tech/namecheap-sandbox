import { useState } from "react";
import { useDomain } from "../context/DomainContext"; // 👈 import

type Props = {
  onSuccess: () => void;
};

export default function BuyDomainForm({ onSuccess }: Props) {
  const { selectedDomain: domain, price } = useDomain(); // 👈 get domain and price

  const [formData, setFormData] = useState({
    firstName: "",
    lastName: "",
    email: "",
    address: "",
    city: "",
    phone: "+91.",
    postalCode: "",
    country: ""
  });

  const [status, setStatus] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleBuy = async () => {
    setStatus("⏳ Processing purchase...");
    try {
      const res = await fetch("http://localhost:8080/buy-domain", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          domain,
          ...formData,
          price: price.base,
          tax: price.tax,
          total: price.total,
          aRecord: "82.25.106.75",
          cName: "indigo-spoonbill-233511.hostingersite.com"
        }),
      });

      const data = await res.json();
      if (res.ok) {
        setStatus(`✅ Domain "${data.domain}" purchased successfully!`);
        onSuccess();
      } else {
        setStatus(`❌ ${data.error || "Purchase failed."}`);
      }
    } catch (err) {
      setStatus("❌ Network or server error.");
    }
  };

  return (
    <div className="bg-white p-6 shadow-md rounded-lg max-w-xl mx-auto mt-6">
      <h2 className="text-xl font-semibold mb-4">
        Enter your details to buy <span className="text-blue-600">{domain}</span>
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <input name="firstName" value={formData.firstName} onChange={handleChange} placeholder="First Name" className="p-2 border rounded" />
        <input name="lastName" value={formData.lastName} onChange={handleChange} placeholder="Last Name" className="p-2 border rounded" />
        <input name="email" value={formData.email} onChange={handleChange} placeholder="Email" className="p-2 border rounded col-span-2" />
        <input name="phone" value={formData.phone} onChange={handleChange} placeholder="Phone" className="p-2 border rounded" />
        <input name="postalCode" value={formData.postalCode} onChange={handleChange} placeholder="Postal Code" className="p-2 border rounded" />
        <input name="address" value={formData.address} onChange={handleChange} placeholder="Address" className="p-2 border rounded col-span-2" />
        <input name="city" value={formData.city} onChange={handleChange} placeholder="City" className="p-2 border rounded col-span-2" />
        <input name="country" value={formData.country} onChange={handleChange} placeholder="Country" className="p-2 border rounded col-span-2" />
      </div>

      <button
        className="bg-green-600 text-white py-2 px-6 rounded hover:bg-green-700 transition"
        onClick={handleBuy}
      >
        Buy for ₹{price.total.toFixed(2)}
      </button>

      {status && <p className="mt-4 text-sm text-gray-700">{status}</p>}
    </div>
  );
}
