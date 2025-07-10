import { useState } from "react";
import DomainSearch from "./components/DomainSearch";
import DomainDetails from "./components/DomainDetails";
import BuyDomainForm from "./components/BuyDomainForm";
import DomainList from "./components/Domainlist";
import EditDNSModal from "./components/EditDNSModal";

export default function App() {
  const [domain, setDomain] = useState("");
  const [step, setStep] = useState(1); // 1: Search, 2: Details, 3: Form, 4: List

  return (
    <div className="min-h-screen w-full bg-gray-100 p-6">
      <div className="max-w-4xl mx-auto grid grid-cols-1 lg:flex-cols-2 gap-6">
        {/* Left: Purchase Flow */}
        <div className="bg-white rounded-2xl shadow-lg p-6">
          <h1 className="text-2xl font-bold text-blue-700 mb-2">Sandbox Domain Reseller</h1>
          <p className="text-gray-600 text-sm mb-6">Search, View, and Purchase your domain</p>

          {step === 1 && (
            <DomainSearch
              onAvailable={(domain) => {
                setDomain(domain);
                setStep(2);
              }}
            />
          )}

          {step === 2 && (
            <DomainDetails
              domain={domain}
              onBuyClick={() => setStep(3)}
            />
          )}

          {step === 3 && (
            <BuyDomainForm
              domain={domain}
              onSuccess={() => {
                alert("✅ Domain bought!");
                setStep(1); // Go back to search
                window.dispatchEvent(new Event("refreshDomains")); // Notify list
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
