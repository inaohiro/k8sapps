export function AppPage({
  content,
  handleReissueToken
}: {
  content: React.ReactNode,
  handleReissueToken: () => void
}) {

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex justify-end p-4 bg-white shadow">
        <button
          onClick={handleReissueToken}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        >
          トークン再発行
        </button>
      </div>
      <div className="max-w-3xl mx-auto mt-8">
        {content}
      </div>
    </div>
  );
}
