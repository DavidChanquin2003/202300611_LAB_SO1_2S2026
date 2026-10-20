use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct WarReport {
    country: String,
    warplanes_in_air: u32,
    warships_in_water: u32,
    timestamp: String,
}

#[post("/grpc-{carnet}")]
async fn receive_report(
    path: web::Path<String>,
    report: web::Json<WarReport>,
) -> impl Responder {
    let carnet = path.into_inner();
    println!("Reporte recibido de Locust en /grpc-{}: {:?}", carnet, report);

    // Enviar el JSON al Cliente Go (que está en Kubernetes como go-client-svc)
    let client = reqwest::Client::new();
    let res = client.post("http://go-client-svc:80/send")
        .json(&report.into_inner())
        .send()
        .await;

    match res {
        Ok(response) => {
            // Si Go responde bien, devolvemos su respuesta al usuario
            let body = response.text().await.unwrap_or_default();
            HttpResponse::Ok().content_type("application/json").body(body)
        },
        Err(e) => {
            println!("Error conectando con Go Client: {}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "status": "error",
                "message": "Fallo al comunicar con Go Client"
            }))
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Iniciando API Rust en el puerto 8080...");
    HttpServer::new(|| {
        App::new().service(receive_report)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}