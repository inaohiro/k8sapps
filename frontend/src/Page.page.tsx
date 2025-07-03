import React, { useState } from "react";

export function Page({ children }: { children: React.ReactNode }) {
  return <div className="max-w-3xl mx-auto mt-8">{children}</div>;
}
