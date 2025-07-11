import { useState } from "react";
import DomainSearch from "./components/DomainSearch";
import DomainDetails from "./components/DomainDetails";
import BuyDomainForm from "./components/BuyDomainForm";
import DomainList from "./components/Domainlist";
import EditDNSModal from "./components/EditDNSModal";

export default function App() {
  const [domain, setDomain] = useState("");
  const [step, setStep] = useState(1); // 1: Search, 2: Details, 3: Form
  const [price, setPrice] = useState({ base: 0, tax: 0, total: 0 });

  const handleAvailable = async (domain: string) => {
    setDomain(domain);
    setStep(2);

    try {
      const res = await fetch(`http://localhost:8080/domain-price?domain=${domain}`);
      const data = await res.json();
      setPrice(data);
    } catch (err) {
      console.error("❌ Failed to fetch domain price:", err);
      setPrice({ base: 1000, tax: 300, total: 1300 }); // fallback
    }
  };

  return (
    <div className="min-h-screen w-full bg-gray-100 p-6">
      <div className="max-w-4xl mx-auto grid grid-cols-1 lg:flex-cols-2 gap-6">
        {/* Left: Purchase Flow */}
        <div className="bg-white rounded-2xl shadow-lg p-6">
          <h1 className="text-2xl font-bold text-blue-700 mb-2">Sandbox Domain Reseller</h1>
          <p className="text-gray-600 text-sm mb-6">Search, View, and Purchase your domain</p>

          {step === 1 && (
            <DomainSearch onAvailable={handleAvailable} />
          )}

          {step === 2 && (
            <DomainDetails
              domain={domain}
              price={price}
              onBuyClick={() => setStep(3)}
            />
          )}

          {step === 3 && (
            <BuyDomainForm
              domain={domain}
              price={price}
              onSuccess={() => {
                alert("✅ Domain bought!");
                setStep(1);
                window.dispatchEvent(new Event("refreshDomains"));
              }}
            />
          )}
        </div>

        {/* Right: Purchased Domains List */}
        <div className="bg-white rounded-2xl shadow-lg p-8">
          <h2 className="text-xl font-semibold mb-4 text-green-700">Your Purchased Domains</h2>
          <DomainList />
        </div>
      </div>

      {/* Global DNS Modal */}
      <EditDNSModal />
    </div>
  );
}
