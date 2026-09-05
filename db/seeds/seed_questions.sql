-- Seed de 50 preguntas para entorno de desarrollo
-- Idempotente: usa UUIDs deterministas + ON CONFLICT DO NOTHING
-- Ejecución: docker exec -i dev-postgres psql -U devuser -d dev_test_api < db/seeds/seed_questions.sql
-- O vía: make db-seed
-- Nota: content/explanation usan markdown (soportado por CodeContent.vue):
--   `código inline`, bloques ```lang ... ``` y **negrita**.

BEGIN;

-- Usuario semilla (dueño de las preguntas). Solo para FK; source=ai_generated las hace visibles a todos.
INSERT INTO users (id, email, password_hash, is_admin, created_at, updated_at)
VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'seed@dev.local', '$2a$10$seedDummyHashSeedDummyHashSeedDum', false, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Preguntas (50)
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000001', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace la palabra clave `defer` en Go?

```go
defer file.Close()
fmt.Println("hola")
```', 'La sentencia **defer** programa una llamada para ejecutarse justo antes de que la función retorne, en orden **LIFO**. Es ideal para liberar recursos como `file.Close()` o `mu.Unlock()`.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000002', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'En Go, ¿qué diferencia hay entre un **array** y un **slice**?

```go
var a [3]int          // array
var s []int = a[:]    // slice
```', 'Un **array** tiene tamaño fijo y forma parte de su tipo (`[3]int`); un **slice** es una vista dinámica sobre un array subyacente con **puntero, longitud y capacidad** (`len`/`cap`).', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000003', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué sucede si envías a un **canal sin buffer** sin receptor listo?

```go
ch := make(chan int)
ch <- 42 // ¿qué ocurre aquí?
```', 'Un canal sin buffer es un **rendezvous síncrono**: el emisor se **bloquea** hasta que otro goroutine ejecute `<-ch`.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000004', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', '¿Cuáles de estas afirmaciones sobre **interfaces** en Go son correctas? (elige todas las válidas)', 'En Go las interfaces se implementan **implícitamente**: un tipo satisface una interfaz si posee todos sus métodos. La interfaz vacía `interface{}` puede contener **cualquier valor**.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000005', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la función en Go que retorna la suma de dos enteros.

```go
func Suma(a, b int) int {
    // completa aquí
}
```', 'Basta con retornar `a + b`. En Go el **retorno nombrado** no es obligatorio; un `return a+b` es suficiente.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000006', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué aporta **TypeScript** sobre JavaScript?', 'Añade **tipado estático opcional** verificado en compilación que se borra al emitir JavaScript. No reemplaza a `V8` ni compila a binario.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000007', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace el genérico `Array<T>` en TypeScript?

```ts
const nums: Array<number> = [1, 2, 3];
```', 'Parametriza el tipo de los elementos del array con `T`, permitiendo reutilizar lógica manteniendo **seguridad de tipos**.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000008', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', '¿Qué **utility types** de TypeScript transforman propiedades? (elige todas las válidas)

```ts
type A = Partial<User>
type B = ReturnType<typeof fn>
```', '`Partial`, `Required`, `Pick` y `Omit` transforman mapeos de propiedades. `ReturnType` extrae el **retorno** de una función.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000009', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la función genérica en TypeScript que retorna el primer elemento de un array.

```ts
function primero<T>(arr: T[]): T | undefined {
  // completa aquí
}
```', 'Usa `T` para tipar el array y el retorno. Retorna `arr[0]` o `undefined` si está vacío.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000010', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es un **closure** en JavaScript?

```js
function makeCounter() {
  let n = 0;
  return () => ++n;
}
```', 'Un **closure** es una función que conserva acceso al **ámbito léxico** donde fue creada, aunque se ejecute fuera de él.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000011', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Cómo funciona el **event loop** en JavaScript?

```js
console.log("1");
setTimeout(() => console.log("2"), 0);
console.log("3");
```', 'El **event loop** toma tareas de la cola y las ejecuta cuando el **call stack** está vacío, permitiendo concurrencia no bloqueante.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000012', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre **Promesas** en JavaScript, ¿cuáles son correctas?

```js
fetch("/api").then(r => r.json()).catch(console.error);
```', '`then`/`catch`/`finally` encadenan; `Promise.all` espera a todas o falla rápido; `async`/`await` es **azúcar sintáctico** sobre promesas.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000013', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la función en JavaScript que filtra los números **pares** de un array.

```js
function pares(arr) {
  // retorna solo los pares
}
```', 'Usa `Array.filter` con `n % 2 === 0`.

```js
return arr.filter(n => n % 2 === 0);
```', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000014', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace una **list comprehension** en Python?

```python
[x*x for x in range(5) if x % 2 == 0]
```', 'Crea una nueva lista aplicando **expresión y filtro** en una sola línea de forma declarativa y eficiente.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000015', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es un **decorador** en Python?

```python
@mi_decorador
def saludar(): ...
```', 'Un decorador es un **callable** que recibe una función y retorna otra envolviéndola, usando la sintaxis `@decorador`.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000016', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es el **GIL** en CPython?', 'El **Global Interpreter Lock** permite que solo un hilo ejecute **bytecode Python** a la vez; afecta al **CPU-bound threading**.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000017', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la función en Python que retorna el **factorial** de `n` (`n >= 0`).

```python
def factorial(n):
    # completa aquí
    pass
```', 'Itera o usa recursión; recuerda que `factorial(0) == 1`.

```python
res = 1
for i in range(2, n+1): res *= i
return res
```', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000018', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace `useState` en React?

```jsx
const [count, setCount] = useState(0);
```', 'Declara **estado local** en un componente funcional; retorna `[valor, setter]` y dispara **re-render** al cambiar.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000019', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'En React, ¿para qué sirve el **array de dependencias** de `useEffect`?

```js
useEffect(() => { fetchData(); }, [userId]);
```', 'Controla cuándo se re-ejecuta el efecto: vacío `[]` = solo montaje, con valores = cuando cambian, **omitido** = cada render.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000020', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', '¿Cuáles son **reglas de los Hooks** en React? (elige todas las válidas)', 'Solo llamar hooks en el **top-level** del componente o custom hook, y solo desde **funciones React**.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000021', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Cuándo usar `useMemo` en React?

```jsx
const total = useMemo(() => items.reduce((a,b) => a+b, 0), [items]);
```', 'Memoriza un **valor computado costoso** entre renders; no es garantía de caché permanente.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000022', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace `v-model` en Vue?

```vue
<input v-model="nombre" />
```', 'Crea **two-way binding** entre un input y una variable reactiva (azúcar sobre `:modelValue` + `@update:modelValue`).', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000023', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué aporta la **Composition API** frente a Options API en Vue 3?

```ts
setup() { const count = ref(0); return { count }; }
```', 'Permite agrupar lógica por **feature** con `setup()`, mejor reutilización y **tipado con TypeScript**.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000024', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre el **ciclo de vida** en Vue 3, ¿cuáles son correctas?

```ts
onMounted(() => { /* ... */ })
```', '`onMounted` tras montar el DOM, `onUpdated` tras actualizar, `onUnmounted` al desmontar; no existe `onBeforeRender`.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000025', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace `computed` en Vue?

```ts
const doble = computed(() => count.value * 2);
```', 'Propiedad **derivada cacheada** que se recalcula solo cuando sus dependencias reactivas cambian.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000026', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace la cláusula `WHERE` en SQL?

```sql
SELECT * FROM users WHERE active = true;
```', 'Filtra **filas** antes de agrupar/ordenar; sin `WHERE` se devuelven todas las filas.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000027', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Diferencia entre `INNER JOIN` y `LEFT JOIN`?

```sql
SELECT * FROM a LEFT JOIN b ON a.id = b.a_id;
```', '`INNER` solo filas con **match** en ambas tablas; `LEFT` conserva todas las de la izquierda aunque no haya match (`NULL`s).', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000028', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Para qué sirve un índice **B-Tree** en SQL?

```sql
CREATE INDEX idx_users_email ON users(email);
```', 'Acelera búsquedas por **igualdad/rango** y ordenamientos a costa de escrituras y espacio.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000029', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre transacciones **ACID**, ¿cuáles son correctas?', '**Atomicity** todo o nada, **Consistency** invariantes, **Isolation** entre transacciones, **Durability** tras `commit`.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000030', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué ventaja ofrece `JSONB` en PostgreSQL frente a `JSON`?

```sql
SELECT data->>''nombre'' FROM docs WHERE data @> ''{"activo": true}'';
```', '`JSONB` almacena **binario parseado**, permite índices `GIN` y operadores eficientes; `JSON` guarda texto original.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000031', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es **MVCC** en PostgreSQL?', '**Multi-Version Concurrency Control** mantiene versiones de filas para **lecturas no bloqueantes** y aislamiento.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000032', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Diferencia entre `SERIAL` y `UUID` como clave primaria en PostgreSQL?

```sql
id SERIAL PRIMARY KEY
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```', '`SERIAL` es entero **autoincremental secuencial**; `UUID` es globalmente único, mejor para distribuidos pero más grande.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000033', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué define un `Dockerfile`?

```dockerfile
FROM node:20-alpine
COPY . /app
RUN npm ci
CMD ["node", "server.js"]
```', 'Receta de **capas** para construir una imagen: `FROM`, `RUN`, `COPY`, `CMD`, etc.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000034', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Para qué sirve un **build multi-stage** en Docker?

```dockerfile
FROM golang:1.22 AS builder
RUN go build -o app .
FROM alpine
COPY --from=builder /app /app
```', 'Usa etapas intermedias para **compilar** y copia solo artefactos a la imagen final, reduciendo tamaño y superficie.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000035', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre **capas de imágenes Docker**, ¿cuáles son correctas?', 'Cada instrucción crea una **capa cacheable**; el orden afecta invalidación de caché; las capas permiten **compartir base** entre imágenes.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000036', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace `chmod 755 archivo` en Linux?

```bash
chmod 755 script.sh
ls -l # -rwxr-xr-x
```', '`755` = **rwx** para owner, **r-x** para group y others (`7=111`, `5=101`).', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000037', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué hace `grep -r "TODO" .` en Linux?

```bash
grep -r "TODO" .
```', 'Busca **recursivamente** la cadena `TODO` en todos los archivos bajo el directorio actual.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000038', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué método HTTP es **idempotente y seguro** para obtener un recurso?

```http
GET /api/users/42 HTTP/1.1
```', '`GET` es **seguro** (no muta) e **idempotente** (múltiples llamadas mismo efecto). `POST` no lo es.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000039', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué código HTTP indica **creación exitosa** de un recurso?

```http
HTTP/1.1 201 Created
Location: /api/users/42
```', '`201 Created` indica recurso creado; `200` es OK genérico, `204` No Content sin cuerpo.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000040', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Diferencia entre `query` y `mutation` en GraphQL?

```graphql
query { user(id: 1) { name } }
mutation { createUser(name: "Ana") { id } }
```', '`query` **lee** datos, `mutation` los **modifica**; ambas se validan contra el schema.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000041', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es el problema **N+1** en GraphQL y cómo se mitiga?

```js
// 1 query para usuarios + N queries para sus posts
users.map(u => fetchPosts(u.id))
```', 'Resolver cada nodo con query adicional genera **N+1 consultas**; se mitiga con **DataLoader/batch**.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000042', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es el **Event Loop** en Node.js?', 'Bucle que gestiona **I/O asíncrono** y callbacks sobre un solo hilo principal con `libuv`.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000043', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué establece el Principio de **Responsabilidad Única** (SRP) de SOLID?', 'Una clase/módulo debe tener **una sola razón para cambiar**, es decir, una sola responsabilidad.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000044', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es el **polimorfismo** en OOP?

```ts
interface Sonido { sonar(): void; }
class Perro implements Sonido { sonar() { console.log("guau"); } }
```', 'Capacidad de tratar objetos de distintas clases vía una **interfaz común**, con comportamiento específico por subtipo.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000045', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Cuál es la complejidad promedio de búsqueda en una **tabla hash** bien dimensionada?

```python
tabla["clave"]  # acceso
```', '`O(1)` promedio; `O(n)` en peor caso por colisiones.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000046', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Diferencia entre **mock** y **stub** en testing?

```ts
stub.returns(42);
expect(mock.calledOnce).toBe(true);
```', '`Stub` retorna respuestas fijas; `mock` además **verifica interacciones/expectativas**.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000047', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Cómo mitigas **XSS almacenado** en una app web?

```js
// ❌ peligroso
el.innerHTML = userInput;
// ✅ seguro
el.textContent = userInput;
```', 'Escapa/sanitiza **salida**, usa **CSP**, y valida entrada; no basta solo con validar en cliente.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000048', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Cuál es la diferencia entre **búsqueda lineal** y **binaria**?

```python
# lineal: O(n)  |  binaria: O(log n) sobre array ordenado
```', 'Lineal `O(n)` recorre todo; binaria `O(log n)` requiere **array ordenado** y divide por la mitad.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000049', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Para qué se usa **Redis** comúnmente?

```bash
SET sesion:abc "{\"user\": 42}" EX 3600
GET sesion:abc
```', 'Almacén **clave-valor en memoria** para **caché**, sesiones, colas y pub/sub con alta velocidad.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000050', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', '¿Qué es un **componente** en Angular?

```ts
@Component({ selector: "app-card", template: "<h1>{{title}}</h1>" })
export class CardComponent { @Input() title = ""; }
```', 'Unidad básica de UI con **template**, estilos y lógica, gestionada por el framework con **DI** y change detection.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Opciones (solo para single_choice / multiple_choice)
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'Ejecuta la función al inicio de la siguiente iteración', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', 'Programa una llamada para ejecutarse al retornar la función (**LIFO**)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', 'Cancela la ejecución de la goroutine actual', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000001', 'Convierte la función en asíncrona', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000002', 'No hay diferencia, son alias', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000002', 'El **array** es de tamaño fijo y el **slice** es dinámico con `puntero/longitud/capacidad`', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000002', 'El slice solo almacena enteros', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000002', 'El array se almacena en el heap y el slice en el stack', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000003', 'El valor se descarta silenciosamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000010', '10000000-0000-0000-0000-000000000003', 'El emisor se **bloquea** hasta que haya un receptor (`rendezvous`)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000003', 'Se genera un `panic` automáticamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000012', '10000000-0000-0000-0000-000000000003', 'El canal crea un buffer temporal', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000004', 'La implementación es **implícita**, no se declara `implements`', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000014', '10000000-0000-0000-0000-000000000004', 'Un tipo debe declarar explícitamente qué interfaces implementa', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000015', '10000000-0000-0000-0000-000000000004', 'Una interfaz vacía `interface{}` puede contener **cualquier valor**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000016', '10000000-0000-0000-0000-000000000004', 'Las interfaces solo pueden contener métodos privados', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000017', '10000000-0000-0000-0000-000000000006', 'Ejecución más rápida en el navegador', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000018', '10000000-0000-0000-0000-000000000006', 'Tipado **estático opcional** verificado en compilación', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000019', '10000000-0000-0000-0000-000000000006', 'Reemplaza al motor `V8`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000020', '10000000-0000-0000-0000-000000000006', 'Convierte JS en lenguaje compilado a binario', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000021', '10000000-0000-0000-0000-000000000007', 'Define un array que solo acepta `strings`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000022', '10000000-0000-0000-0000-000000000007', 'Permite que una estructura trabaje con **cualquier tipo** manteniendo el tipado (`<T>`)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000023', '10000000-0000-0000-0000-000000000007', 'Convierte el array en inmutable', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000024', '10000000-0000-0000-0000-000000000007', 'Es un alias de `any[]`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000025', '10000000-0000-0000-0000-000000000008', '`Partial<T>` hace todas las propiedades **opcionales**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000026', '10000000-0000-0000-0000-000000000008', '`Pick<T,K>` selecciona un **subconjunto** de claves', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000027', '10000000-0000-0000-0000-000000000008', '`ReturnType<T>` transforma propiedades en opcionales', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000028', '10000000-0000-0000-0000-000000000008', '`Omit<T,K>` **excluye** claves del tipo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000029', '10000000-0000-0000-0000-000000000010', 'Un bucle que se cierra automáticamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000030', '10000000-0000-0000-0000-000000000010', 'Una función que **recuerda el ámbito** donde fue creada', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000031', '10000000-0000-0000-0000-000000000010', 'Un método para cerrar el navegador', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000032', '10000000-0000-0000-0000-000000000010', 'Una forma de declarar clases', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000033', '10000000-0000-0000-0000-000000000011', 'Ejecuta todo en paralelo con hilos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000034', '10000000-0000-0000-0000-000000000011', 'Procesa la cola de tareas cuando el **call stack** está vacío', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000035', '10000000-0000-0000-0000-000000000011', 'Compila JS a bytecode antes de ejecutar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000036', '10000000-0000-0000-0000-000000000011', 'Solo gestiona eventos del DOM', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000037', '10000000-0000-0000-0000-000000000012', '`then` encadena el valor resuelto y `catch` captura rechazos', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000038', '10000000-0000-0000-0000-000000000012', '`Promise.all` espera a todas o rechaza al primer fallo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000039', '10000000-0000-0000-0000-000000000012', '`await` solo funciona dentro de funciones `async` (o top-level en módulos)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000040', '10000000-0000-0000-0000-000000000012', 'Las promesas se ejecutan en un hilo separado del event loop', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000041', '10000000-0000-0000-0000-000000000014', 'Ordena la lista in-place', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000042', '10000000-0000-0000-0000-000000000014', 'Crea una nueva lista aplicando **expresión/filtro** en una línea', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000043', '10000000-0000-0000-0000-000000000014', 'Convierte la lista en tupla', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000044', '10000000-0000-0000-0000-000000000014', 'Elimina duplicados automáticamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000045', '10000000-0000-0000-0000-000000000015', 'Un comentario especial que documenta la función', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000046', '10000000-0000-0000-0000-000000000015', 'Una función que **envuelve otra** para extender su comportamiento (`@decorador`)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000047', '10000000-0000-0000-0000-000000000015', 'Una palabra clave para declarar clases abstractas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000048', '10000000-0000-0000-0000-000000000015', 'Un módulo de la stdlib para logging', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000049', '10000000-0000-0000-0000-000000000016', 'Un recolector de basura generacional', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000050', '10000000-0000-0000-0000-000000000016', 'Un lock global que permite solo un hilo ejecutando **bytecode** a la vez', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000051', '10000000-0000-0000-0000-000000000016', 'Una optimización JIT', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000052', '10000000-0000-0000-0000-000000000016', 'Un protocolo de serialización', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000053', '10000000-0000-0000-0000-000000000018', 'Gestiona rutas de la aplicación', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000054', '10000000-0000-0000-0000-000000000018', 'Declara **estado local** que dispara re-render al actualizarse (`[valor, setter]`)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000055', '10000000-0000-0000-0000-000000000018', 'Hace peticiones HTTP', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000056', '10000000-0000-0000-0000-000000000018', 'Memoriza componentes hijos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000057', '10000000-0000-0000-0000-000000000019', 'Define el estado inicial del efecto', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000058', '10000000-0000-0000-0000-000000000019', 'Controla **cuándo** se re-ejecuta el efecto según valores observados', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000059', '10000000-0000-0000-0000-000000000019', 'Lista los componentes hijos a renderizar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000060', '10000000-0000-0000-0000-000000000019', 'Especifica los props que el efecto puede mutar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000061', '10000000-0000-0000-0000-000000000020', 'Llamar hooks solo en el **nivel superior**, no dentro de `loops`/`condicionales`', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000062', '10000000-0000-0000-0000-000000000020', 'Llamar hooks solo desde **componentes React** o custom hooks', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000063', '10000000-0000-0000-0000-000000000020', 'Los hooks pueden llamarse desde cualquier función JS', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000064', '10000000-0000-0000-0000-000000000020', 'El orden de llamada de hooks puede variar entre renders', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000065', '10000000-0000-0000-0000-000000000021', 'Para evitar re-renders de componentes hijos siempre', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000066', '10000000-0000-0000-0000-000000000021', 'Para memorizar un **cálculo costoso** entre renders', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000067', '10000000-0000-0000-0000-000000000021', 'Para reemplazar a `useEffect`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000068', '10000000-0000-0000-0000-000000000021', 'Para hacer el componente asíncrono', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000069', '10000000-0000-0000-0000-000000000022', 'Renderiza una lista de elementos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000070', '10000000-0000-0000-0000-000000000022', 'Sincroniza un input con una variable reactiva (**two-way binding**)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000071', '10000000-0000-0000-0000-000000000022', 'Define una ruta del router', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000072', '10000000-0000-0000-0000-000000000022', 'Registra un componente global', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000073', '10000000-0000-0000-0000-000000000023', 'Elimina la reactividad del sistema', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000074', '10000000-0000-0000-0000-000000000023', 'Permite organizar lógica por **feature** con mejor reutilización y tipado', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000075', '10000000-0000-0000-0000-000000000023', 'Solo funciona con Vue 2', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000076', '10000000-0000-0000-0000-000000000023', 'Obliga a usar clases', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000077', '10000000-0000-0000-0000-000000000024', '`onMounted` se ejecuta **tras montar** el componente en el DOM', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000078', '10000000-0000-0000-0000-000000000024', '`onUpdated` se ejecuta tras **actualizar** el DOM reactivo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000079', '10000000-0000-0000-0000-000000000024', '`onBeforeRender` es un hook oficial de Vue 3', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000080', '10000000-0000-0000-0000-000000000024', '`onUnmounted` se ejecuta al **desmontar** el componente', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000081', '10000000-0000-0000-0000-000000000025', 'Ejecuta un efecto secundario sin retorno', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000082', '10000000-0000-0000-0000-000000000025', 'Deriva un valor **cacheado** que se recalcula solo si cambian sus dependencias', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000083', '10000000-0000-0000-0000-000000000025', 'Registra un watcher manual', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000084', '10000000-0000-0000-0000-000000000025', 'Define un componente asíncrono', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000085', '10000000-0000-0000-0000-000000000026', 'Ordena los resultados', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000086', '10000000-0000-0000-0000-000000000026', 'Filtra filas que cumplen una **condición** (`WHERE`)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000087', '10000000-0000-0000-0000-000000000026', 'Agrupa filas por columna', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000088', '10000000-0000-0000-0000-000000000026', 'Define claves foráneas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000089', '10000000-0000-0000-0000-000000000027', 'No hay diferencia', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000090', '10000000-0000-0000-0000-000000000027', '`INNER` retorna solo **coincidencias**; `LEFT` conserva todas las filas de la izquierda', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000091', '10000000-0000-0000-0000-000000000027', '`LEFT` es más rápido siempre', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000092', '10000000-0000-0000-0000-000000000027', '`INNER` solo funciona con claves primarias', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000093', '10000000-0000-0000-0000-000000000028', 'Encripta la columna', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000094', '10000000-0000-0000-0000-000000000028', 'Acelera **búsquedas y ordenamientos** por esa columna', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000095', '10000000-0000-0000-0000-000000000028', 'Garantiza unicidad sin declarar `UNIQUE`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000096', '10000000-0000-0000-0000-000000000028', 'Convierte la tabla en temporal', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000097', '10000000-0000-0000-0000-000000000029', '**Atomicity**: todo o nada', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000098', '10000000-0000-0000-0000-000000000029', '**Consistency**: se preservan invariantes', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000099', '10000000-0000-0000-0000-000000000029', '**Isolation**: transacciones no se interfieren según nivel', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000100', '10000000-0000-0000-0000-000000000029', '**Durability**: `commit` persiste aunque falle el sistema', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000101', '10000000-0000-0000-0000-000000000030', '`JSONB` es solo un alias de `TEXT`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000102', '10000000-0000-0000-0000-000000000030', '`JSONB` almacena **binario indexable** y más eficiente para consultas', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000103', '10000000-0000-0000-0000-000000000030', '`JSON` no permite anidar objetos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000104', '10000000-0000-0000-0000-000000000030', '`JSONB` no soporta índices', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000105', '10000000-0000-0000-0000-000000000031', 'Un tipo de índice para texto', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000106', '10000000-0000-0000-0000-000000000031', 'Control de concurrencia por **versiones** que evita bloqueos de lectura', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000107', '10000000-0000-0000-0000-000000000031', 'Un protocolo de replicación síncrona', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000108', '10000000-0000-0000-0000-000000000031', 'Un motor de full-text search', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000109', '10000000-0000-0000-0000-000000000032', 'Son idénticos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000110', '10000000-0000-0000-0000-000000000032', '`SERIAL` es entero secuencial; `UUID` es globalmente único y no secuencial', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000111', '10000000-0000-0000-0000-000000000032', '`UUID` solo funciona en MySQL', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000112', '10000000-0000-0000-0000-000000000032', '`SERIAL` no puede ser clave primaria', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000113', '10000000-0000-0000-0000-000000000033', 'Un contenedor en ejecución', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000114', '10000000-0000-0000-0000-000000000033', 'Instrucciones para construir una **imagen Docker**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000115', '10000000-0000-0000-0000-000000000033', 'Un volumen persistente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000116', '10000000-0000-0000-0000-000000000033', 'La configuración del daemon', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000117', '10000000-0000-0000-0000-000000000034', 'Para ejecutar varios contenedores a la vez', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000118', '10000000-0000-0000-0000-000000000034', 'Para compilar en una etapa y copiar artefactos a una **imagen final más pequeña**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000119', '10000000-0000-0000-0000-000000000034', 'Para hacer backup de imágenes', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000120', '10000000-0000-0000-0000-000000000034', 'Para paralelizar `docker pull`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000121', '10000000-0000-0000-0000-000000000035', 'Cada instrucción del `Dockerfile` crea una **capa cacheable**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000122', '10000000-0000-0000-0000-000000000035', 'El orden de instrucciones afecta la **invalidación de caché**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000123', '10000000-0000-0000-0000-000000000035', 'Juntar todo en un solo `RUN` siempre reduce el tamaño final', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000124', '10000000-0000-0000-0000-000000000035', 'Las capas permiten **compartir base** entre imágenes', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000125', '10000000-0000-0000-0000-000000000036', 'Solo lectura para todos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000126', '10000000-0000-0000-0000-000000000036', '`rwx` para owner, `r-x` para grupo y otros', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000127', '10000000-0000-0000-0000-000000000036', 'Borra el archivo', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000128', '10000000-0000-0000-0000-000000000036', 'Cambia el propietario a `root`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000129', '10000000-0000-0000-0000-000000000037', 'Reemplaza `TODO` por vacío', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000130', '10000000-0000-0000-0000-000000000037', 'Busca **recursivamente** `TODO` en el árbol de archivos', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000131', '10000000-0000-0000-0000-000000000037', 'Lista archivos modificados por `git`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000132', '10000000-0000-0000-0000-000000000037', 'Muestra el manual de `grep`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000133', '10000000-0000-0000-0000-000000000038', '`POST`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000134', '10000000-0000-0000-0000-000000000038', '`GET` — **seguro** e **idempotente**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000135', '10000000-0000-0000-0000-000000000038', '`PATCH` sin validación', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000136', '10000000-0000-0000-0000-000000000038', '`CONNECT`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000137', '10000000-0000-0000-0000-000000000039', '`200 OK`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000138', '10000000-0000-0000-0000-000000000039', '`201 Created`', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000139', '10000000-0000-0000-0000-000000000039', '`301 Moved Permanently`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000140', '10000000-0000-0000-0000-000000000039', '`400 Bad Request`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000141', '10000000-0000-0000-0000-000000000040', 'No hay diferencia', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000142', '10000000-0000-0000-0000-000000000040', '`query` **lee** datos, `mutation` los **modifica**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000143', '10000000-0000-0000-0000-000000000040', '`query` es REST y `mutation` es GraphQL', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000144', '10000000-0000-0000-0000-000000000040', '`mutation` es solo para suscripciones', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000145', '10000000-0000-0000-0000-000000000041', 'Falta de paginación', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000146', '10000000-0000-0000-0000-000000000041', 'Múltiples queries por cada nodo hijo; se mitiga con **batch/DataLoader**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000147', '10000000-0000-0000-0000-000000000041', 'Error de tipado en el schema', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000148', '10000000-0000-0000-0000-000000000041', 'Exceso de mutaciones anidadas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000149', '10000000-0000-0000-0000-000000000042', 'Un hilo que compila TypeScript', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000150', '10000000-0000-0000-0000-000000000042', 'Bucle que gestiona **operaciones asíncronas** sobre un solo hilo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000151', '10000000-0000-0000-0000-000000000042', 'Un balanceador de carga', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000152', '10000000-0000-0000-0000-000000000042', 'Un recolector de basura', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000153', '10000000-0000-0000-0000-000000000043', 'Una clase debe tener muchos motivos para cambiar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000154', '10000000-0000-0000-0000-000000000043', 'Una clase debe tener **una sola razón** para cambiar', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000155', '10000000-0000-0000-0000-000000000043', 'Prohíbe usar herencia', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000156', '10000000-0000-0000-0000-000000000043', 'Obliga a usar singletons', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000157', '10000000-0000-0000-0000-000000000044', 'Ocultar datos internos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000158', '10000000-0000-0000-0000-000000000044', 'Tratar objetos distintos mediante una **interfaz común** con comportamientos específicos', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000159', '10000000-0000-0000-0000-000000000044', 'Duplicar código en subclases', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000160', '10000000-0000-0000-0000-000000000044', 'Convertir objetos a JSON', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000161', '10000000-0000-0000-0000-000000000045', '`O(n)`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000162', '10000000-0000-0000-0000-000000000045', '`O(1)` promedio', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000163', '10000000-0000-0000-0000-000000000045', '`O(log n)`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000164', '10000000-0000-0000-0000-000000000045', '`O(n log n)`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000165', '10000000-0000-0000-0000-000000000046', 'Son sinónimos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000166', '10000000-0000-0000-0000-000000000046', '`Stub` retorna datos fijos; `mock` **verifica interacciones** y expectativas', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000167', '10000000-0000-0000-0000-000000000046', '`Mock` solo funciona en integración', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000168', '10000000-0000-0000-0000-000000000046', '`Stub` requiere Docker', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000169', '10000000-0000-0000-0000-000000000047', 'Solo validando en el frontend', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000170', '10000000-0000-0000-0000-000000000047', 'Escapando **salida**, sanitizando HTML y usando **Content-Security-Policy**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000171', '10000000-0000-0000-0000-000000000047', 'Deshabilitando JavaScript en el navegador', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000172', '10000000-0000-0000-0000-000000000047', 'Usando solo HTTP sin TLS', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000173', '10000000-0000-0000-0000-000000000048', 'Son iguales', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000174', '10000000-0000-0000-0000-000000000048', 'Lineal `O(n)` recorre todo; binaria `O(log n)` requiere **array ordenado**', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000175', '10000000-0000-0000-0000-000000000048', 'Binaria solo funciona en listas enlazadas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000176', '10000000-0000-0000-0000-000000000048', 'Lineal requiere array ordenado', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000177', '10000000-0000-0000-0000-000000000049', 'Solo como base de datos relacional', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000178', '10000000-0000-0000-0000-000000000049', '**Caché en memoria**, sesiones y colas de alta velocidad', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000179', '10000000-0000-0000-0000-000000000049', 'Compilador de Go', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000180', '10000000-0000-0000-0000-000000000049', 'Balanceador de carga', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000181', '10000000-0000-0000-0000-000000000050', 'Un archivo CSS', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000182', '10000000-0000-0000-0000-000000000050', 'Unidad de UI con **template**, lógica y estilos gestionada por Angular', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000183', '10000000-0000-0000-0000-000000000050', 'Una tabla de base de datos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000184', '10000000-0000-0000-0000-000000000050', 'Un hook de React', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;

-- Code challenges (para code_completion)
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', 'package main

func Suma(a, b int) int {
    // completa aquí
}
', '5', 'go', '[]', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000009', 'function primero<T>(arr: T[]): T | undefined {
  // completa aquí
}
', '1', 'typescript', '[]', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000013', 'function pares(arr) {
  // retorna solo los pares
}
', '[2,4]', 'javascript', '[]', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000017', '10000000-0000-0000-0000-000000000017', 'def factorial(n):
    # completa aquí
    pass
', '120', 'python', '[]', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;

-- Relaciones pregunta <-> topic (por slug, UUIDs de topics son aleatorios)
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000001', id, NOW() FROM topics WHERE slug='go' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000002', id, NOW() FROM topics WHERE slug='go' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000003', id, NOW() FROM topics WHERE slug='go' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000004', id, NOW() FROM topics WHERE slug='go' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000005', id, NOW() FROM topics WHERE slug='go' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000006', id, NOW() FROM topics WHERE slug='typescript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000007', id, NOW() FROM topics WHERE slug='typescript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000008', id, NOW() FROM topics WHERE slug='typescript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000009', id, NOW() FROM topics WHERE slug='typescript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000010', id, NOW() FROM topics WHERE slug='javascript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000011', id, NOW() FROM topics WHERE slug='javascript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000012', id, NOW() FROM topics WHERE slug='javascript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000013', id, NOW() FROM topics WHERE slug='javascript' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000014', id, NOW() FROM topics WHERE slug='python' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000015', id, NOW() FROM topics WHERE slug='python' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000016', id, NOW() FROM topics WHERE slug='python' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000017', id, NOW() FROM topics WHERE slug='python' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000018', id, NOW() FROM topics WHERE slug='react' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000019', id, NOW() FROM topics WHERE slug='react' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000020', id, NOW() FROM topics WHERE slug='react' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000021', id, NOW() FROM topics WHERE slug='react' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000022', id, NOW() FROM topics WHERE slug='vue' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000023', id, NOW() FROM topics WHERE slug='vue' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000024', id, NOW() FROM topics WHERE slug='vue' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000025', id, NOW() FROM topics WHERE slug='vue' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000026', id, NOW() FROM topics WHERE slug='sql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000027', id, NOW() FROM topics WHERE slug='sql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000028', id, NOW() FROM topics WHERE slug='sql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000029', id, NOW() FROM topics WHERE slug='sql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000030', id, NOW() FROM topics WHERE slug='postgresql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000031', id, NOW() FROM topics WHERE slug='postgresql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000032', id, NOW() FROM topics WHERE slug='postgresql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000033', id, NOW() FROM topics WHERE slug='docker' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000034', id, NOW() FROM topics WHERE slug='docker' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000035', id, NOW() FROM topics WHERE slug='docker' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000036', id, NOW() FROM topics WHERE slug='linux' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000037', id, NOW() FROM topics WHERE slug='linux' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000038', id, NOW() FROM topics WHERE slug='rest' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000039', id, NOW() FROM topics WHERE slug='rest' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000040', id, NOW() FROM topics WHERE slug='graphql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000041', id, NOW() FROM topics WHERE slug='graphql' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000042', id, NOW() FROM topics WHERE slug='nodejs' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000043', id, NOW() FROM topics WHERE slug='solid' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000044', id, NOW() FROM topics WHERE slug='oop' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000045', id, NOW() FROM topics WHERE slug='estructuras-datos' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000046', id, NOW() FROM topics WHERE slug='testing' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000047', id, NOW() FROM topics WHERE slug='seguridad' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000048', id, NOW() FROM topics WHERE slug='algoritmos' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000049', id, NOW() FROM topics WHERE slug='redis' ON CONFLICT DO NOTHING;
INSERT INTO question_topics (question_id, topic_id, created_at) SELECT '10000000-0000-0000-0000-000000000050', id, NOW() FROM topics WHERE slug='angular' ON CONFLICT DO NOTHING;

COMMIT;
