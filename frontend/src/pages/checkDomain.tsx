import { useState } from "react";
import { useNavigate } from "react-router-dom";
import DomainSearch from "../pages/DomainSearch";
import DomainDetails from "../pages/DomainDetails";

export default function CheckDomain() {
  const [domain, setDomain] = useState("");
  const [price, setPrice] = useState({ base: 0, tax: 0, total: 0 });
  const navigate = useNavigate();

  const handleAvailable = async (domain: string) => {
    setDomain(domain);
    try {
      const res = await fetch(`http://localhost:8080/domain-price?domain=${domain}`);
      const data = await res.json();
      setPrice(data);
    } catch (err) {
      console.error("❌ Failed to fetch price:", err);
      setPrice({ base: 1000, tax: 300, total: 1300 });
    }
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow">
      <h2 className="text-2xl font-bold text-blue-700 mb-4">Search & Buy a Domain</h2>
      <DomainSearch onAvailable={handleAvailable} />
      {domain && (
        <DomainDetails
          domain={domain}
          price={price}
          onBuyClick={() => navigate("/buy", { state: { domain, price } })}
        />
      )}
    </div>
  );
}
