import http from "k6/http";

export const options = {
    vus: 30,
    duration: "30s",
};

export default function () {
    http.get("http://127.0.0.1:3000/random-by-time-range");
}
