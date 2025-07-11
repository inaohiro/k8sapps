import { PageState } from "./store/store";

export function Sidebar({
  setPage,
}: {
  setPage: (t: PageState) => void;
}) {
  return (

  <div className="fixed left-0 top-0 w-64 h-full bg-gray-900 p-4">
    <ul className="mt-4">
      <li className="mb-1 group active">
        <a onClick={() => setPage({ type: "deployments-list" })}
        href="#" className="flex items-center py-2 px-4 text-gray-300 hover:bg-gray-950 rounded-md group-[.active]:text-white">
          <i className="ri-home-2-line mr-3 text-lg"></i>
          <span className="text-sm">Deployments</span>
        </a>
      </li>
      <li className="mb-1 group active">
        <a onClick={() => setPage({ type: "services-list" })}
        href="#" className="flex items-center py-2 px-4 text-gray-300 hover:bg-gray-950 rounded-md group-[.active]:text-white">
          <i className="ri-home-2-line mr-3 text-lg"></i>
          <span className="text-sm">Services</span>
        </a>
      </li>
      <li className="mb-1 group active">
        <a onClick={() => setPage({ type: "pods-list" })}
        href="#" className="flex items-center py-2 px-4 text-gray-300 hover:bg-gray-950 rounded-md group-[.active]:text-white">
          <i className="ri-home-2-line mr-3 text-lg"></i>
          <span className="text-sm">Pods</span>
        </a>
      </li>
    </ul>
  </div>
  );
}
