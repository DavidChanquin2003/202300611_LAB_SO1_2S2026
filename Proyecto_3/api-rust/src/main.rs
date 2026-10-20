use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

// Definimos la estructura del JSON que enviará Locust
#[derive(Debug, Serialize, Deserialize)]
struct WarReport {
    country: String,
    warplanes_in_air: u32,
    warships_in_water: u32,
    timestamp: String,
}

// Endpoint dinámico para recibir el reporte
#[post("/grpc-{carnet}")]
async fn receive_report(
    path: web::Path<String>,
    report: web::Json<WarReport>,
) -> impl Responder {
    let carnet = path.into_inner();
    
    // Imprimimos el log para confirmar recepción
    println!("Reporte recibido en /grpc-{}: {:?}", carnet, report);

    // Retornamos 200 OK
    HttpResponse::Ok().json(serde_json::json!({
        "status": "success",
        "message": format!("Reporte procesado en ruta de carnet {}", carnet)
    }))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Iniciando API Rust en el puerto 8080...");
    HttpServer::new(|| {
        App::new()
            .service(receive_report)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}