package database

import (
	"log"

	"github.com/felipe/dev-test-api/internal/models"
	"gorm.io/gorm"
)

func seedDefaultTopics(db *gorm.DB) {
	type seedTopic struct {
		Slug     string
		Name     string
		Category string
	}

	topics := []seedTopic{
		// --- lenguajes (11) ---
		{"go", "Go", "lenguajes"},
		{"typescript", "TypeScript", "lenguajes"},
		{"javascript", "JavaScript", "lenguajes"},
		{"python", "Python", "lenguajes"},
		{"java", "Java", "lenguajes"},
		{"rust", "Rust", "lenguajes"},
		{"csharp", "C#", "lenguajes"},
		{"php", "PHP", "lenguajes"},
		{"ruby", "Ruby", "lenguajes"},
		{"kotlin", "Kotlin", "lenguajes"},
		{"swift", "Swift", "lenguajes"},
		{"mojo", "Mojo", "lenguajes"},
		{"zig", "Zig", "lenguajes"},
		// --- frontend (11) ---
		{"react", "React", "frontend"},
		{"nextjs", "Next.js", "frontend"},
		{"vue", "Vue", "frontend"},
		{"angular", "Angular", "frontend"},
		{"svelte", "Svelte", "frontend"},
		{"css", "CSS", "frontend"},
		{"html", "HTML", "frontend"},
		{"tailwind", "Tailwind CSS", "frontend"},
		{"astro", "Astro", "frontend"},
		{"nuxt", "Nuxt", "frontend"},
		{"vite", "Vite", "frontend"},
		// --- backend (12) ---
		{"nodejs", "Node.js", "backend"},
		{"express", "Express", "backend"},
		{"django", "Django", "backend"},
		{"spring-boot", "Spring Boot", "backend"},
		{"gin", "Gin", "backend"},
		{"laravel", "Laravel", "backend"},
		{"fastapi", "FastAPI", "backend"},
		{"graphql", "GraphQL", "backend"},
		{"rest", "REST", "backend"},
		{"nestjs", "NestJS", "backend"},
		{"grpc", "gRPC", "backend"},
		{"bun", "Bun", "backend"},
		// --- base-datos (6) ---
		{"postgresql", "PostgreSQL", "base-datos"},
		{"sql", "SQL", "base-datos"},
		{"mongodb", "MongoDB", "base-datos"},
		{"redis", "Redis", "base-datos"},
		{"prisma", "Prisma", "base-datos"},
		{"gorm", "GORM", "base-datos"},
		// --- datos / ingeniería de datos (10) ---
		{"ingenieria-datos", "Ingeniería de Datos", "datos"},
		{"analisis-datos", "Análisis de Datos", "datos"},
		{"etl-elt", "ETL / ELT", "datos"},
		{"airflow", "Apache Airflow", "datos"},
		{"spark", "Apache Spark", "datos"},
		{"kafka", "Apache Kafka", "datos"},
		{"snowflake", "Snowflake", "datos"},
		{"databricks", "Databricks", "datos"},
		{"dbt", "dbt", "datos"},
		{"data-warehousing", "Data Warehousing", "datos"},
		// --- devops (10) ---
		{"docker", "Docker", "devops"},
		{"kubernetes", "Kubernetes", "devops"},
		{"ci-cd", "CI/CD", "devops"},
		{"aws", "AWS", "devops"},
		{"terraform", "Terraform", "devops"},
		{"linux", "Linux", "devops"},
		{"github-actions", "GitHub Actions", "devops"},
		{"argocd", "ArgoCD", "devops"},
		{"opentelemetry", "OpenTelemetry", "devops"},
		{"ansible", "Ansible", "devops"},
		// --- cloud (6) ---
		{"gcp", "GCP", "cloud"},
		{"azure", "Azure", "cloud"},
		{"serverless", "Serverless", "cloud"},
		{"vercel", "Vercel", "cloud"},
		{"cloud-native", "Cloud Native", "cloud"},
		{"observabilidad", "Observabilidad", "cloud"},
		// --- arquitectura (8) ---
		{"microservicios", "Microservicios", "arquitectura"},
		{"ddd", "Domain-Driven Design", "arquitectura"},
		{"solid", "Principios SOLID", "arquitectura"},
		{"patrones-diseno", "Patrones de Diseño", "arquitectura"},
		{"clean-architecture", "Clean Architecture", "arquitectura"},
		{"event-driven", "Event-Driven Architecture", "arquitectura"},
		{"hexagonal", "Arquitectura Hexagonal", "arquitectura"},
		{"cqrs", "CQRS", "arquitectura"},
		// --- conceptos (7) ---
		{"algoritmos", "Algoritmos", "conceptos"},
		{"estructuras-datos", "Estructuras de Datos", "conceptos"},
		{"system-design", "System Design", "conceptos"},
		{"testing", "Testing", "conceptos"},
		{"seguridad", "Seguridad", "conceptos"},
		{"oop", "Programación Orientada a Objetos", "conceptos"},
		{"funcional", "Programación Funcional", "conceptos"},
		// --- testing (6) ---
		{"testing-unitario", "Testing Unitario", "testing"},
		{"testing-e2e", "Testing E2E", "testing"},
		{"playwright", "Playwright", "testing"},
		{"jest-vitest", "Jest / Vitest", "testing"},
		{"tdd", "TDD", "testing"},
		{"cypress", "Cypress", "testing"},
		// --- seguridad (8) ---
		{"devsecops", "DevSecOps", "seguridad"},
		{"owasp", "OWASP Top 10", "seguridad"},
		{"pentesting", "Pentesting", "seguridad"},
		{"oauth-oidc", "OAuth / OIDC", "seguridad"},
		{"zero-trust", "Zero Trust", "seguridad"},
		{"criptografia", "Criptografía", "seguridad"},
		{"appsec", "AppSec", "seguridad"},
		{"ethical-hacking", "Ethical Hacking", "seguridad"},
		// --- ia (13) ---
		{"llm", "LLM", "ia"},
		{"rag", "RAG", "ia"},
		{"ai-agents", "AI Agents", "ia"},
		{"prompt-engineering", "Prompt Engineering", "ia"},
		{"vector-db", "Vector DB", "ia"},
		{"mlops", "MLOps", "ia"},
		{"genai", "Generative AI", "ia"},
		{"nlp", "NLP", "ia"},
		{"computer-vision", "Computer Vision", "ia"},
		{"mcp", "Model Context Protocol", "ia"},
		{"langchain", "LangChain", "ia"},
		{"embeddings", "Embeddings", "ia"},
		{"machine-learning", "Machine Learning", "ia"},
		// --- movil (5) ---
		{"flutter", "Flutter", "movil"},
		{"react-native", "React Native", "movil"},
		{"android", "Android", "movil"},
		{"ios", "iOS", "movil"},
		{"kotlin-multiplatform", "Kotlin Multiplatform", "movil"},
	}

	for _, t := range topics {
		db.Where(map[string]interface{}{
			"slug":       t.Slug,
			"created_by": nil,
		}).Attrs(models.Topic{
			Name:     t.Name,
			Category: t.Category,
			IsSystem: true,
		}).FirstOrCreate(&models.Topic{})
	}

	log.Printf("🌱 Default topics seeded: %d available", len(topics))
}
