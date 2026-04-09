package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/ewertonfrnc/social-network/internal/store"
)

var usernames = []string{
	"ana", "bruno", "carla", "daniel", "elis", "felipe", "gabriela", "henrique",
	"isabela", "joao", "karina", "lucas", "mariana", "natalia", "otavio", "patricia",
	"quezia", "rafael", "samara", "thiago", "ursula", "vinicius", "wanda", "xavier",
	"yasmin", "zeca", "aline", "bia", "caio", "diego", "eduardo", "fernanda",
	"gustavo", "helena", "igor", "juliana", "kaua", "leticia", "mateus", "nina",
	"oliveira", "paulo", "quirino", "rodrigo", "silvia", "tatiane", "ubirajara", "valeria",
	"willian", "ximena", "yuri", "zilda",
}

var titles = []string{
	"O Poder do Habito", "Abracando o Minimalismo", "Dicas de Alimentacao Saudavel",
	"Viajar com Orcamento", "Meditacao Mindfulness", "Aumente Sua Produtividade",
	"Montando Seu Home Office", "Detox Digital", "Nocoes de Jardinagem",
	"Projetos DIY para Casa", "Yoga para Iniciantes", "Vida Sustentavel",
	"Dominando a Gestao do Tempo", "Explorando a Natureza", "Receitas Simples para Cozinhar",
	"Treino em Casa", "Dicas de Financas Pessoais", "Escrita Criativa",
	"Consciencia sobre Saude Mental", "Aprendendo Novas Habilidades",
}

var contents = []string{
	"Neste post, vamos explorar como criar bons habitos que permanecem e transformam sua vida.",
	"Descubra os beneficios de um estilo de vida minimalista e como reduzir excessos em casa e na mente.",
	"Aprenda dicas praticas para se alimentar bem gastando pouco, sem abrir mao do sabor.",
	"Viajar nao precisa ser caro. Aqui vao dicas para conhecer o mundo com um orcamento menor.",
	"A meditacao mindfulness pode reduzir o estresse e melhorar seu bem-estar mental. Veja como comecar.",
	"Aumente sua produtividade com estrategias simples e eficazes.",
	"Monte o home office ideal para melhorar eficiencia e conforto no trabalho remoto.",
	"Um detox digital pode ajudar voce a se reconectar com o mundo real e cuidar melhor da saude mental.",
	"Comece sua jornada na jardinagem com estas dicas basicas para iniciantes.",
	"Transforme sua casa com projetos DIY divertidos e faceis de fazer.",
	"Yoga e uma otima forma de manter o corpo ativo e flexivel. Veja posturas para iniciantes.",
	"Viver de forma sustentavel e bom para voce e para o planeta. Aprenda escolhas mais ecologicas.",
	"Domine a gestao do tempo com estas dicas e faca mais em menos tempo.",
	"A natureza tem muito a oferecer. Descubra os beneficios de passar mais tempo ao ar livre.",
	"Prepare refeicoes deliciosas com receitas simples e rapidas.",
	"Mantenha a forma sem sair de casa com treinos eficientes.",
	"Assuma o controle do seu dinheiro com dicas praticas de financas pessoais.",
	"Solte sua criatividade com propostas inspiradoras de escrita e exercicios.",
	"Saude mental e tao importante quanto saude fisica. Aprenda a cuidar melhor da sua mente.",
	"Aprender novas habilidades pode ser divertido e recompensador. Veja ideias para comecar.",
}

var tags = []string{
	"Autodesenvolvimento", "Minimalismo", "Saude", "Viagem", "Mindfulness",
	"Produtividade", "Home Office", "Detox Digital", "Jardinagem", "DIY",
	"Yoga", "Sustentabilidade", "Gestao do Tempo", "Natureza", "Culinaria",
	"Treino", "Financas Pessoais", "Escrita", "Saude Mental", "Aprendizado",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	_ = db

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("Erro ao criar usuario:", err)
			_ = tx.Rollback()
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Erro ao criar post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Erro ao criar comentario:", err)
			return
		}
	}

	log.Println("Seed concluido")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range num {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := range num {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)
	for i := range comments {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]

		comments[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: contents[rand.Intn(len(contents))],
		}
	}

	return comments
}
