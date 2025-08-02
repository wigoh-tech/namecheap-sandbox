import { useState } from "react";
import { useNavigate } from "react-router-dom";
import DomainSearch from "../pages/DomainSearch";
import DomainDetails from "../pages/DomainDetails";
import { useDomain } from "../context/DomainContext"; // 👈 import

export default function CheckDomain() {
  const navigate = useNavigate();
  const { setSelectedDomain, setPrice } = useDomain(); // 👈 use context

  const handleAvailable = async (domain: string) => {
    setSelectedDomain(domain);
    try {
      const res = await fetch(`http://localhost:8080/domain-price?domain=${domain}`);
      const data = await res.json();
      setPrice(data);
    } catch (err) {
      console.error("❌ Failed to fetch price:", err);
      setPrice({ base: 1000, tax: 300, total: 1300 }); // fallback
    }
    navigate("/buy"); // 👈 go to buy page
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow">
      <h2 className="text-2xl font-bold text-blue-700 mb-4">Search & Buy a Domain</h2>
      <DomainSearch onAvailable={handleAvailable} />
    </div>
  );
}
