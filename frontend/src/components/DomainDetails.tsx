import { useEffect, useState } from "react";

type Props = {
  domain: string;
  onBuyClick: () => void;
};

export default function DomainDetails({ domain, onBuyClick }: Props) {
  const [price, setPrice] = useState({ base: 0, tax: 0, total: 0 });

  useEffect(() => {
    const fetchPrice = async () => {
      try {
        const res = await fetch(`http://localhost:8080/domain-price?domain=${domain}`);
        const data = await res.json();
        setPrice(data);
      } catch (err) {
        console.error("Failed to fetch price", err);
        setPrice({ base: 1000, tax: 300, total: 1300 }); // fallback
      }
    };

    if (domain) {
      fetchPrice();
    }
  }, [domain]);

  return (
    <div className="bg-white border p-4 rounded-lg shadow space-y-3">
      <h2 className="text-xl font-semibold text-green-600">Domain Available 🎉</h2>
      <p className="text-gray-800">You can purchase <strong>{domain}</strong></p>
      <p className="text-gray-700">Price: ₹{price.total.toFixed(2)}</p>
      <button
        className="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-6 rounded transition"
        onClick={onBuyClick}
      >
        Proceed to Buy
      </button>
    </div>
  );
}
