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
		{"react", "React", "frontend"},
		{"nextjs", "Next.js", "frontend"},
		{"vue", "Vue", "frontend"},
		{"angular", "Angular", "frontend"},
		{"svelte", "Svelte", "frontend"},
		{"css", "CSS", "frontend"},
		{"html", "HTML", "frontend"},
		{"tailwind", "Tailwind CSS", "frontend"},
		{"nodejs", "Node.js", "backend"},
		{"express", "Express", "backend"},
		{"django", "Django", "backend"},
		{"spring-boot", "Spring Boot", "backend"},
		{"gin", "Gin", "backend"},
		{"laravel", "Laravel", "backend"},
		{"fastapi", "FastAPI", "backend"},
		{"graphql", "GraphQL", "backend"},
		{"rest", "REST", "backend"},
		{"docker", "Docker", "devops"},
		{"kubernetes", "Kubernetes", "devops"},
		{"ci-cd", "CI/CD", "devops"},
		{"aws", "AWS", "devops"},
		{"terraform", "Terraform", "devops"},
		{"linux", "Linux", "devops"},
		{"microservicios", "Microservicios", "arquitectura"},
		{"ddd", "Domain-Driven Design", "arquitectura"},
		{"solid", "Principios SOLID", "arquitectura"},
		{"patrones-diseno", "Patrones de Diseño", "arquitectura"},
		{"clean-architecture", "Clean Architecture", "arquitectura"},
		{"event-driven", "Event-Driven Architecture", "arquitectura"},
		{"postgresql", "PostgreSQL", "base-datos"},
		{"sql", "SQL", "base-datos"},
		{"mongodb", "MongoDB", "base-datos"},
		{"redis", "Redis", "base-datos"},
		{"prisma", "Prisma", "base-datos"},
		{"gorm", "GORM", "base-datos"},
		{"algoritmos", "Algoritmos", "conceptos"},
		{"estructuras-datos", "Estructuras de Datos", "conceptos"},
		{"system-design", "System Design", "conceptos"},
		{"testing", "Testing", "conceptos"},
		{"seguridad", "Seguridad", "conceptos"},
		{"oop", "Programación Orientada a Objetos", "conceptos"},
		{"funcional", "Programación Funcional", "conceptos"},
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

	log.Printf("Default topics seeded: %d available", len(topics))
}
