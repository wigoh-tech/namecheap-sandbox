type Props = {
  domain: string;
  onBuyClick: () => void;
};

export default function DomainDetails({ domain, onBuyClick }: Props) {
  const base = 1000;
  const tax = base * 0.3;
  const total = base + tax;

  return (
    <div className="bg-white border p-4 rounded-lg shadow space-y-3">
      <h2 className="text-xl font-semibold text-green-600">Domain Available 🎉</h2>
      <p className="text-gray-800">You can purchase <strong>{domain}</strong></p>
      <p className="text-gray-700">Price: ₹{total.toFixed(2)} (₹{base} + 30%)</p>
      <button
        className="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-6 rounded transition"
        onClick={onBuyClick}
      >
        Proceed to Buy
      </button>
    </div>
  );
}
