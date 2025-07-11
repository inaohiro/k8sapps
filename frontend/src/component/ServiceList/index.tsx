import { useSetAtom } from "jotai";
import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { pageAtom, setPageAtom } from "../../store/store";

interface Service {
  id: string;
  name: string;
  type: string;
  clusterIP: string;
  ports: ServicePort[];
}
type ServicePort = {
  port: string;
  targetPort: string;
  protocol: null | string;
};

export function ServiceList() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

  useEffect(() => {
    fetch(`/api/services`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch services");
        return res.json();
      })
      .then((data) => {
        setServices(data);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Services</h1>
      <button
        onClick={() => {
          setPage({ type: "services-create" });
        }}
        style={{ marginBottom: "1em" }}
      >
        + 新規 Service 作成
      </button>

      <div className="relative overflow-x-auto">
        <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
          <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" className="px-6 py-3">
                Name
              </th>
              <th scope="col" className="px-6 py-3">
                Type
              </th>
              <th scope="col" className="px-6 py-3">
                IP
              </th>
              <th scope="col" className="px-6 py-3">
                Ports
              </th>
              <th scope="col" className="px-6 py-3">
                Action
              </th>
            </tr>
          </thead>
          <tbody>
            {services.map((svc) => (
              <tr key={svc.id} className="bg-white border-b dark:bg-gray-800 dark:border-gray-700 border-gray-200">
                <th scope="row" className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      setPage({ type: "services-detail", id: svc.name });
                    }}
                    style={{ fontWeight: "bold", cursor: "pointer", color: "#1976d2", textDecoration: "underline" }}
                  >
                    {svc.name}
                  </a>
                </th>
                <td className="px-6 py-4">{svc.type}</td>
                <td className="px-6 py-4">{svc.clusterIP}</td>
                <td className="px-6 py-4">
                  {svc.ports.map((port) => (
                    <>
                      <span>
                        Port: {port.port}, TargetPort: {port.targetPort}, Protocol: {port.protocol ?? "tcp"}
                      </span>
                      <br />
                    </>
                  ))}
                </td>
                <td className="px-6 py-4">
                  <button
                    style={{
                      marginLeft: 8,
                      color: "#fff",
                      background: "#d32f2f",
                      border: "none",
                      borderRadius: 4,
                      padding: "2px 8px",
                      cursor: "pointer",
                    }}
                    onClick={async (e) => {
                      e.stopPropagation();
                      setLoading(true);
                      setError(null);
                      try {
                        const res = await fetch(`/api/services/${svc.name}`, {
                          method: "DELETE",
                          headers: { Authorization: `Bearer ${token}` },
                        });
                        if (!res.ok) throw new Error("Failed to delete service");
                      } catch (err: any) {
                        setError(err.message);
                      } finally {
                        setLoading(false);
                      }
                    }}
                  >
                    削除
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
