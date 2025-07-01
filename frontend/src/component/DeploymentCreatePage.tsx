import React, { useEffect, useState } from "react";
import { useToken } from "../hooks/useToken";

interface Image {
  name: string;
  tag: string;
}

const DeploymentCreatePage: React.FC = () => {
  const [name, setName] = useState("");
  const [image, setImage] = useState("");
  const [images, setImages] = useState<Image[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const { token } = useToken();

  useEffect(() => {
    fetch("/api/images", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch images");
        return res.json();
      })
      .then((data: Image[]) => {
        setImages(data);
        if (data.length > 0) setImage(`${data[0].name}:${data[0].tag}`);
      })
      .catch((err: Error) => {
        setError(err.message);
      });
  }, [token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    try {
      const res = await fetch(`/api/deployments`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ name, image }),
      });
      if (!res.ok) throw new Error("Failed to create deployment");
      setSuccess(true);
      setName("");
      setImage(images.length > 0 ? `${images[0].name}:${images[0].tag}` : "");
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h1>新規 Deployment 作成</h1>
      <form onSubmit={handleSubmit} style={{ maxWidth: 400 }}>
        <div style={{ marginBottom: 12 }}>
          <label>
            名前
            <br />
            <input value={name} onChange={(e) => setName(e.target.value)} required />
          </label>
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>
            イメージ
            <br />
            <select value={image} onChange={(e) => setImage(e.target.value)} required>
              {images.map((img) => (
                <option key={`${img.name}:${img.tag}`} value={`${img.name}:${img.tag}`}>
                  {img.name}:{img.tag}
                </option>
              ))}
            </select>
          </label>
        </div>
        <button type="submit">作成</button>
      </form>
      {error && <div style={{ color: "red" }}>Error: {error}</div>}
      {success && <div style={{ color: "green" }}>Deployment を作成しました</div>}
    </div>
  );
};

export default DeploymentCreatePage;
