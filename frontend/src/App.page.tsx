import React, { useState } from "react";

export function AppPage({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* <div className="max-w-3xl mx-auto mt-8">{content}</div> */}
      {children}
    </div>
  );
}
