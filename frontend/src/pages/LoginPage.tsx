import { useToast } from "../context/ToastContext";

export default function LoginPage() {
  const { showToast } = useToast();

  return (
    <>
      <button onClick={() => showToast(`Welcome back`, "success")}>
        Show Success Toast
      </button>
    </>
  );
}
