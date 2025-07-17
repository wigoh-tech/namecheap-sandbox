import { useEffect, useState } from "react";

type Props = {
  domain: string;
  price: { base: number; tax: number; total: number };
  onBuyClick: () => void;
};

export default function DomainDetails({ domain, price, onBuyClick }: Props) {
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