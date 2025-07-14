import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { useSetAtom } from "jotai";
import { setPageAtom } from "src/store/store";

type Image  = {
  name: string;
}
type Flavor  = {
  name: string;
}

export function DeploymentCreate() {
  const [name, setName] = useState("");
  const [image, setImage] = useState("");
  const [images, setImages] = useState<Image[]>([]);
  const [flavor, setFlavor] = useState("");
  const [flavors, setFlavors] = useState<Flavor[]>([]);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

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
        if (data.length > 0) setImage(`${data[0].name}`);
      })
      .catch((err: Error) => {
        setError(err.message);
      });
  }, [token]);

  useEffect(() => {
    fetch("/api/flavors", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch flavors");
        return res.json();
      })
      .then((data: Flavor[]) => {
        setFlavors(data);
        if (data.length > 0) setFlavor(`${data[0].name}`);
      })
      .catch((err: Error) => {
        setError(err.message);
      });
  }, [token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      const res = await fetch(`/api/deployments`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ name, image }),
      });
      if (!res.ok) throw new Error("Failed to create deployment");
      setPage({type: "deployments-list"})
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
                <option key={`${img.name}`} value={`${img.name}`}>
                  {img.name}
                </option>
              ))}
            </select>
          </label>
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>
            プラン
            <br />
            <select value={flavor} onChange={(e) => setImage(e.target.value)} required>
              {flavors.map((f) => (
                <option key={`${f.name}`} value={`${f.name}`}>
                  {f.name}
                </option>
              ))}
            </select>
          </label>
        </div>
        <button type="submit">作成</button>
      </form>
      {error && <div style={{ color: "red" }}>Error: {error}</div>}
    </div>
  );
};
