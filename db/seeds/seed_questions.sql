-- Seed de 50 preguntas para entorno de desarrollo
-- Idempotente: usa UUIDs deterministas + ON CONFLICT DO NOTHING
-- Ejecucion: docker exec -i dev-postgres psql -U devuser -d dev_test_api < db/seeds/seed_questions.sql
-- O via: make db-seed

BEGIN;

-- Usuario semilla (dueno de las preguntas). Solo para FK; source=ai_generated las hace visibles a todos.
INSERT INTO users (id, email, password_hash, is_admin, created_at, updated_at)
VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'seed@dev.local', '$2a$10$seedDummyHashSeedDummyHashSeedDum', false, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Preguntas (50)
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000001', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace la palabra clave `defer` en Go?', 'defer programa una llamada para ejecutarse justo antes de que la funcion retorne, en orden LIFO. Es ideal para liberar recursos.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000002', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'En Go, que diferencia hay entre un slice y un array?', 'Un array tiene tamano fijo parte de su tipo; un slice es una vista dinamica sobre un array subyacente (puntero, longitud y capacidad).', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000003', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que sucede si envias a un canal sin buffer sin receptor listo?', 'Un canal sin buffer bloquea al emisor hasta que otro goroutine reciba. Es un rendezvous sincrono.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000004', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Cuales de estas afirmaciones sobre interfaces en Go son correctas? (elige todas las validas)', 'En Go las interfaces se implementan implicitamente; un tipo implementa una interfaz si posee todos sus metodos.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000005', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la funcion en Go que retorna la suma de dos enteros.', 'Basta retornar a+b. En Go el retorno nombrado no es obligatorio.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000006', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que aporta TypeScript sobre JavaScript?', 'Anade tipado estatico opcional que se verifica en compilacion y se borra al emitir JavaScript.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000007', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace el generico `Array<T>` en TypeScript?', 'Parametriza el tipo de los elementos del array, permitiendo reutilizar logica manteniendo seguridad de tipos.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000008', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Que utility types de TypeScript transforman propiedades? (elige todas las validas)', 'Partial, Required, Pick y Omit transforman mapeos de propiedades. ReturnType extrae el retorno de una funcion.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000009', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la funcion generica en TypeScript que retorna el primer elemento de un array.', 'Usa T para tipar el array y el retorno. Retorna arr[0] o undefined si esta vacio.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000010', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es un closure en JavaScript?', 'Un closure es una funcion que conserva acceso al ambito lexico donde fue creada, aunque se ejecute fuera de el.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000011', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Como funciona el event loop en JavaScript?', 'El event loop toma tareas de la cola y las ejecuta cuando el call stack esta vacio, permitiendo concurrencia no bloqueante.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000012', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre Promesas en JavaScript, cuales son correctas?', 'then/catch/finally encadenan; Promise.all espera a todas o falla rapido; async/await es azucar sintactico.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000013', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la funcion en JavaScript que filtra los numeros pares de un array.', 'Usa Array.filter con n % 2 === 0.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000014', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace una list comprehension en Python?', 'Crea una nueva lista aplicando expresion y filtro en una sola linea de forma declarativa y eficiente.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000015', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es un decorador en Python?', 'Un decorador es un callable que recibe una funcion y retorna otra envolviendola, usando la sintaxis @decorador.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000016', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es el GIL en CPython?', 'El Global Interpreter Lock permite que solo un hilo ejecute bytecode Python a la vez; afecta CPU-bound threading.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000017', (SELECT id FROM users WHERE email='seed@dev.local'), 'code_completion', 'Completa la funcion en Python que retorna el factorial de n (n >= 0).', 'Itera o usa recursion; factorial(0)=1.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000018', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace `useState` en React?', 'Declara estado local en un componente funcional; retorna [valor, setter] y dispara re-render al cambiar.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000019', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'En React, para que sirve el array de dependencias de `useEffect`?', 'Controla cuando se re-ejecuta el efecto: vacio = solo montaje, con valores = cuando cambian, omitido = cada render.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000020', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Cuales son reglas de los Hooks en React? (elige todas las validas)', 'Solo llamar hooks en el top-level del componente o custom hook, y solo desde funciones React.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000021', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Cuando usar `useMemo` en React?', 'Memoriza un valor computado costoso entre renders; no es garantia de cache permanente.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000022', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace `v-model` en Vue?', 'Crea two-way binding entre un input y una variable reactiva (azucar sobre :value + @update).', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000023', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que aporta la Composition API frente a Options API en Vue 3?', 'Permite agrupar logica por feature con setup(), mejor reutilizacion y tipado con TypeScript.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000024', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre el ciclo de vida en Vue 3, cuales son correctas?', 'onMounted tras montar DOM, onUpdated tras actualizar, onUnmounted al desmontar; no existe onBeforeRender.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000025', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace `computed` en Vue?', 'Propiedad derivada cacheada que se recalcula solo cuando sus dependencias reactivas cambian.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000026', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace la clausula `WHERE` en SQL?', 'Filtra filas antes de agrupar/ordenar; sin WHERE se devuelven todas las filas.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000027', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Diferencia entre `INNER JOIN` y `LEFT JOIN`?', 'INNER solo filas con match en ambas tablas; LEFT conserva todas las de la izquierda aunque no haya match (NULLs).', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000028', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Para que sirve un indice B-Tree en SQL?', 'Acelera busquedas por igualdad/rango y ordenamientos a costa de escrituras y espacio.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000029', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre transacciones ACID, cuales son correctas?', 'Atomicity todo o nada, Consistency invariantes, Isolation entre transacciones, Durability tras commit.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000030', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que ventaja ofrece `JSONB` en PostgreSQL frente a `JSON`?', 'JSONB almacena binario parseado, permite indices GIN y operadores eficientes; JSON guarda texto original.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000031', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es MVCC en PostgreSQL?', 'Multi-Version Concurrency Control mantiene versiones de filas para lecturas no bloqueantes y aislamiento.', 'advanced', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000032', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Diferencia entre `SERIAL` y `UUID` como clave primaria en PostgreSQL?', 'SERIAL es entero autoincremental secuencial; UUID es globalmente unico, mejor para distribuidos pero mas grande.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000033', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que define un `Dockerfile`?', 'Receta de capas para construir una imagen: FROM, RUN, COPY, CMD, etc.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000034', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Para que sirve un build multi-stage en Docker?', 'Usa etapas intermedias para compilar y copia solo artefactos a la imagen final, reduciendo tamano y superficie.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000035', (SELECT id FROM users WHERE email='seed@dev.local'), 'multiple_choice', 'Sobre capas de imagenes Docker, cuales son correctas?', 'Cada instruccion crea una capa cacheable; el orden afecta invalidacion de cache; menos capas no siempre es mejor.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000036', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace `chmod 755 archivo` en Linux?', 'rwx para owner, r-x para group y others (7=111, 5=101).', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000037', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que hace `grep -r "TODO" .` en Linux?', 'Busca recursivamente la cadena TODO en todos los archivos bajo el directorio actual.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000038', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que metodo HTTP es idempotente y seguro para obtener un recurso?', 'GET es seguro (no muta) e idempotente (multiples llamadas mismo efecto). POST no lo es.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000039', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que codigo HTTP indica creacion exitosa de un recurso?', '201 Created indica recurso creado; 200 es OK generico, 204 No Content sin cuerpo.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000040', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Diferencia entre `query` y `mutation` en GraphQL?', 'query lee datos, mutation los modifica; ambas se validan contra el schema.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000041', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es el problema N+1 en GraphQL y como se mitiga?', 'Resolver cada nodo con query adicional genera N+1 consultas; se mitiga con DataLoader/batch.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000042', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es el Event Loop en Node.js?', 'Bucle que gestiona I/O asincrono y callbacks sobre un solo hilo principal con libuv.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000043', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que establece el Principio de Responsabilidad Unica (SRP) de SOLID?', 'Una clase/modulo debe tener una sola razon para cambiar, es decir, una sola responsabilidad.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000044', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es el polimorfismo en OOP?', 'Capacidad de tratar objetos de distintas clases via una interfaz comun, con comportamiento especifico por subtipo.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000045', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Cual es la complejidad promedio de busqueda en una tabla hash bien dimensionada?', 'O(1) promedio; O(n) en peor caso por colisiones.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000046', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Diferencia entre mock y stub en testing?', 'Stub retorna respuestas fijas; mock ademas verifica interacciones/expectativas.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000047', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Como mitigas XSS almacenado en una app web?', 'Escapa/sanitiza salida, usa CSP, y valida entrada; no basta solo con validar en cliente.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000048', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Cual es la diferencia entre busqueda lineal y binaria?', 'Lineal O(n) recorre todo; binaria O(log n) requiere array ordenado y divide por la mitad.', 'intermediate', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000049', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Para que se usa Redis comunmente?', 'Almacen clave-valor en memoria para cache, sesiones, colas y pub/sub con alta velocidad.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO questions (id, user_id, type, content, explanation, difficulty, is_public, source, created_at, updated_at)
VALUES ('10000000-0000-0000-0000-000000000050', (SELECT id FROM users WHERE email='seed@dev.local'), 'single_choice', 'Que es un componente en Angular?', 'Unidad basica de UI con template, estilos y logica, gestionada por el framework con DI y change detection.', 'beginner', true, 'ai_generated', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Opciones (solo para single_choice / multiple_choice)
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'Ejecuta la funcion al inicio de la siguiente iteracion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', 'Programa una llamada para ejecutarse al retornar la funcion (LIFO)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', 'Cancela la ejecucion de la goroutine actual', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000001', 'Convierte la funcion en asincrona', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000002', 'No hay diferencia, son alias', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000002', 'El array es de tamano fijo y el slice es dinamico con puntero/longitud/capacidad', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000002', 'El slice solo almacena enteros', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000002', 'El array se almacena en el heap y el slice en el stack', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000003', 'El valor se descarta silenciosamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000010', '10000000-0000-0000-0000-000000000003', 'El emisor se bloquea hasta que haya un receptor', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000003', 'Se genera un panic automaticamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000012', '10000000-0000-0000-0000-000000000003', 'El canal crea un buffer temporal', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000004', 'La implementacion es implicita, no se declara `implements`', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000014', '10000000-0000-0000-0000-000000000004', 'Un tipo debe declarar explicitamente que interfaces implementa', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000015', '10000000-0000-0000-0000-000000000004', 'Una interfaz vacia `interface{}` puede contener cualquier valor', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000016', '10000000-0000-0000-0000-000000000004', 'Las interfaces solo pueden contener metodos privados', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000017', '10000000-0000-0000-0000-000000000006', 'Ejecucion mas rapida en el navegador', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000018', '10000000-0000-0000-0000-000000000006', 'Tipado estatico opcional verificado en compilacion', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000019', '10000000-0000-0000-0000-000000000006', 'Reemplaza al motor V8', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000020', '10000000-0000-0000-0000-000000000006', 'Convierte JS en lenguaje compilado a binario', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000021', '10000000-0000-0000-0000-000000000007', 'Define un array que solo acepta strings', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000022', '10000000-0000-0000-0000-000000000007', 'Permite que una estructura trabaje con cualquier tipo manteniendo el tipado', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000023', '10000000-0000-0000-0000-000000000007', 'Convierte el array en inmutable', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000024', '10000000-0000-0000-0000-000000000007', 'Es un alias de `any[]`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000025', '10000000-0000-0000-0000-000000000008', 'Partial<T> hace todas las propiedades opcionales', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000026', '10000000-0000-0000-0000-000000000008', 'Pick<T,K> selecciona un subconjunto de claves', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000027', '10000000-0000-0000-0000-000000000008', 'ReturnType<T> transforma propiedades en opcionales', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000028', '10000000-0000-0000-0000-000000000008', 'Omit<T,K> excluye claves del tipo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000029', '10000000-0000-0000-0000-000000000010', 'Un bucle que se cierra automaticamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000030', '10000000-0000-0000-0000-000000000010', 'Una funcion que recuerda el ambito donde fue creada', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000031', '10000000-0000-0000-0000-000000000010', 'Un metodo para cerrar el navegador', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000032', '10000000-0000-0000-0000-000000000010', 'Una forma de declarar clases', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000033', '10000000-0000-0000-0000-000000000011', 'Ejecuta todo en paralelo con hilos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000034', '10000000-0000-0000-0000-000000000011', 'Procesa la cola de tareas cuando el call stack esta vacio', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000035', '10000000-0000-0000-0000-000000000011', 'Compila JS a bytecode antes de ejecutar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000036', '10000000-0000-0000-0000-000000000011', 'Solo gestiona eventos del DOM', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000037', '10000000-0000-0000-0000-000000000012', '`then` encadena el valor resuelto y `catch` captura rechazos', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000038', '10000000-0000-0000-0000-000000000012', '`Promise.all` espera a todas o rechaza al primer fallo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000039', '10000000-0000-0000-0000-000000000012', '`await` solo funciona dentro de funciones `async` (o top-level en modulos)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000040', '10000000-0000-0000-0000-000000000012', 'Las promesas se ejecutan en un hilo separado del event loop', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000041', '10000000-0000-0000-0000-000000000014', 'Ordena la lista in-place', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000042', '10000000-0000-0000-0000-000000000014', 'Crea una nueva lista aplicando expresion/filtro en una linea', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000043', '10000000-0000-0000-0000-000000000014', 'Convierte la lista en tupla', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000044', '10000000-0000-0000-0000-000000000014', 'Elimina duplicados automaticamente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000045', '10000000-0000-0000-0000-000000000015', 'Un comentario especial que documenta la funcion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000046', '10000000-0000-0000-0000-000000000015', 'Una funcion que envuelve otra para extender su comportamiento', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000047', '10000000-0000-0000-0000-000000000015', 'Una palabra clave para declarar clases abstractas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000048', '10000000-0000-0000-0000-000000000015', 'Un modulo de la stdlib para logging', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000049', '10000000-0000-0000-0000-000000000016', 'Un recolector de basura generacional', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000050', '10000000-0000-0000-0000-000000000016', 'Un lock global que permite solo un hilo ejecutando bytecode a la vez', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000051', '10000000-0000-0000-0000-000000000016', 'Una optimizacion JIT', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000052', '10000000-0000-0000-0000-000000000016', 'Un protocolo de serializacion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000053', '10000000-0000-0000-0000-000000000018', 'Gestiona rutas de la aplicacion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000054', '10000000-0000-0000-0000-000000000018', 'Declara estado local que dispara re-render al actualizarse', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000055', '10000000-0000-0000-0000-000000000018', 'Hace peticiones HTTP', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000056', '10000000-0000-0000-0000-000000000018', 'Memoriza componentes hijos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000057', '10000000-0000-0000-0000-000000000019', 'Define el estado inicial del efecto', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000058', '10000000-0000-0000-0000-000000000019', 'Controla cuando se re-ejecuta el efecto segun valores observados', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000059', '10000000-0000-0000-0000-000000000019', 'Lista los componentes hijos a renderizar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000060', '10000000-0000-0000-0000-000000000019', 'Especifica los props que el efecto puede mutar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000061', '10000000-0000-0000-0000-000000000020', 'Llamar hooks solo en el nivel superior, no dentro de loops/condicionales', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000062', '10000000-0000-0000-0000-000000000020', 'Llamar hooks solo desde componentes React o custom hooks', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000063', '10000000-0000-0000-0000-000000000020', 'Los hooks pueden llamarse desde cualquier funcion JS', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000064', '10000000-0000-0000-0000-000000000020', 'El orden de llamada de hooks puede variar entre renders', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000065', '10000000-0000-0000-0000-000000000021', 'Para evitar re-renders de componentes hijos siempre', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000066', '10000000-0000-0000-0000-000000000021', 'Para memorizar un calculo costoso entre renders', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000067', '10000000-0000-0000-0000-000000000021', 'Para reemplazar a `useEffect`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000068', '10000000-0000-0000-0000-000000000021', 'Para hacer el componente asincrono', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000069', '10000000-0000-0000-0000-000000000022', 'Renderiza una lista de elementos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000070', '10000000-0000-0000-0000-000000000022', 'Sincroniza un input con una variable reactiva (two-way binding)', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000071', '10000000-0000-0000-0000-000000000022', 'Define una ruta del router', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000072', '10000000-0000-0000-0000-000000000022', 'Registra un componente global', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000073', '10000000-0000-0000-0000-000000000023', 'Elimina la reactividad del sistema', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000074', '10000000-0000-0000-0000-000000000023', 'Permite organizar logica por feature con mejor reutilizacion y tipado', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000075', '10000000-0000-0000-0000-000000000023', 'Solo funciona con Vue 2', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000076', '10000000-0000-0000-0000-000000000023', 'Obliga a usar clases', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000077', '10000000-0000-0000-0000-000000000024', '`onMounted` se ejecuta tras montar el componente en el DOM', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000078', '10000000-0000-0000-0000-000000000024', '`onUpdated` se ejecuta tras actualizar el DOM reactivo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000079', '10000000-0000-0000-0000-000000000024', '`onBeforeRender` es un hook oficial de Vue 3', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000080', '10000000-0000-0000-0000-000000000024', '`onUnmounted` se ejecuta al desmontar el componente', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000081', '10000000-0000-0000-0000-000000000025', 'Ejecuta un efecto secundario sin retorno', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000082', '10000000-0000-0000-0000-000000000025', 'Deriva un valor cacheado que se recalcula solo si cambian sus dependencias', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000083', '10000000-0000-0000-0000-000000000025', 'Registra un watcher manual', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000084', '10000000-0000-0000-0000-000000000025', 'Define un componente asincrono', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000085', '10000000-0000-0000-0000-000000000026', 'Ordena los resultados', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000086', '10000000-0000-0000-0000-000000000026', 'Filtra filas que cumplen una condicion', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000087', '10000000-0000-0000-0000-000000000026', 'Agrupa filas por columna', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000088', '10000000-0000-0000-0000-000000000026', 'Define claves foraneas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000089', '10000000-0000-0000-0000-000000000027', 'No hay diferencia', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000090', '10000000-0000-0000-0000-000000000027', 'INNER retorna solo coincidencias; LEFT conserva todas las filas de la izquierda', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000091', '10000000-0000-0000-0000-000000000027', 'LEFT es mas rapido siempre', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000092', '10000000-0000-0000-0000-000000000027', 'INNER solo funciona con claves primarias', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000093', '10000000-0000-0000-0000-000000000028', 'Encripta la columna', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000094', '10000000-0000-0000-0000-000000000028', 'Acelera busquedas y ordenamientos por esa columna', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000095', '10000000-0000-0000-0000-000000000028', 'Garantiza unicidad sin declarar UNIQUE', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000096', '10000000-0000-0000-0000-000000000028', 'Convierte la tabla en temporal', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000097', '10000000-0000-0000-0000-000000000029', 'Atomicidad: todo o nada', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000098', '10000000-0000-0000-0000-000000000029', 'Consistencia: se preservan invariantes', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000099', '10000000-0000-0000-0000-000000000029', 'Aislamiento: transacciones no se interfieren segun nivel', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000100', '10000000-0000-0000-0000-000000000029', 'Durabilidad: commit persiste aunque falle el sistema', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000101', '10000000-0000-0000-0000-000000000030', 'JSONB es solo un alias de TEXT', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000102', '10000000-0000-0000-0000-000000000030', 'JSONB almacena binario indexable y mas eficiente para consultas', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000103', '10000000-0000-0000-0000-000000000030', 'JSON no permite anidar objetos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000104', '10000000-0000-0000-0000-000000000030', 'JSONB no soporta indices', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000105', '10000000-0000-0000-0000-000000000031', 'Un tipo de indice para texto', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000106', '10000000-0000-0000-0000-000000000031', 'Control de concurrencia por versiones que evita bloqueos de lectura', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000107', '10000000-0000-0000-0000-000000000031', 'Un protocolo de replicacion sincrona', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000108', '10000000-0000-0000-0000-000000000031', 'Un motor de full-text search', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000109', '10000000-0000-0000-0000-000000000032', 'Son identicos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000110', '10000000-0000-0000-0000-000000000032', 'SERIAL es entero secuencial; UUID es globalmente unico y no secuencial', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000111', '10000000-0000-0000-0000-000000000032', 'UUID solo funciona en MySQL', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000112', '10000000-0000-0000-0000-000000000032', 'SERIAL no puede ser clave primaria', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000113', '10000000-0000-0000-0000-000000000033', 'Un contenedor en ejecucion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000114', '10000000-0000-0000-0000-000000000033', 'Instrucciones para construir una imagen Docker', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000115', '10000000-0000-0000-0000-000000000033', 'Un volumen persistente', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000116', '10000000-0000-0000-0000-000000000033', 'La configuracion del daemon', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000117', '10000000-0000-0000-0000-000000000034', 'Para ejecutar varios contenedores a la vez', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000118', '10000000-0000-0000-0000-000000000034', 'Para compilar en una etapa y copiar artefactos a una imagen final mas pequena', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000119', '10000000-0000-0000-0000-000000000034', 'Para hacer backup de imagenes', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000120', '10000000-0000-0000-0000-000000000034', 'Para paralelizar `docker pull`', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000121', '10000000-0000-0000-0000-000000000035', 'Cada instruccion del Dockerfile crea una capa cacheable', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000122', '10000000-0000-0000-0000-000000000035', 'El orden de instrucciones afecta la invalidacion de cache', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000123', '10000000-0000-0000-0000-000000000035', 'Juntar todo en un solo RUN siempre reduce el tamano final', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000124', '10000000-0000-0000-0000-000000000035', 'Las capas permiten compartir base entre imagenes', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000125', '10000000-0000-0000-0000-000000000036', 'Solo lectura para todos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000126', '10000000-0000-0000-0000-000000000036', 'rwx para owner, r-x para grupo y otros', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000127', '10000000-0000-0000-0000-000000000036', 'Borra el archivo', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000128', '10000000-0000-0000-0000-000000000036', 'Cambia el propietario a root', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000129', '10000000-0000-0000-0000-000000000037', 'Reemplaza TODO por vacio', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000130', '10000000-0000-0000-0000-000000000037', 'Busca recursivamente TODO en el arbol de archivos', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000131', '10000000-0000-0000-0000-000000000037', 'Lista archivos modificados por git', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000132', '10000000-0000-0000-0000-000000000037', 'Muestra el manual de grep', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000133', '10000000-0000-0000-0000-000000000038', 'POST', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000134', '10000000-0000-0000-0000-000000000038', 'GET', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000135', '10000000-0000-0000-0000-000000000038', 'PATCH sin validacion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000136', '10000000-0000-0000-0000-000000000038', 'CONNECT', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000137', '10000000-0000-0000-0000-000000000039', '200 OK', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000138', '10000000-0000-0000-0000-000000000039', '201 Created', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000139', '10000000-0000-0000-0000-000000000039', '301 Moved Permanently', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000140', '10000000-0000-0000-0000-000000000039', '400 Bad Request', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000141', '10000000-0000-0000-0000-000000000040', 'No hay diferencia', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000142', '10000000-0000-0000-0000-000000000040', 'query lee datos, mutation los modifica', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000143', '10000000-0000-0000-0000-000000000040', 'query es REST y mutation es GraphQL', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000144', '10000000-0000-0000-0000-000000000040', 'mutation es solo para suscripciones', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000145', '10000000-0000-0000-0000-000000000041', 'Falta de paginacion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000146', '10000000-0000-0000-0000-000000000041', 'Multiples queries por cada nodo hijo; se mitiga con batch/DataLoader', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000147', '10000000-0000-0000-0000-000000000041', 'Error de tipado en el schema', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000148', '10000000-0000-0000-0000-000000000041', 'Exceso de mutaciones anidadas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000149', '10000000-0000-0000-0000-000000000042', 'Un hilo que compila TypeScript', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000150', '10000000-0000-0000-0000-000000000042', 'Bucle que gestiona operaciones asincronas sobre un solo hilo', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000151', '10000000-0000-0000-0000-000000000042', 'Un balanceador de carga', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000152', '10000000-0000-0000-0000-000000000042', 'Un recolector de basura', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000153', '10000000-0000-0000-0000-000000000043', 'Una clase debe tener muchos motivos para cambiar', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000154', '10000000-0000-0000-0000-000000000043', 'Una clase debe tener una sola razon para cambiar', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000155', '10000000-0000-0000-0000-000000000043', 'Prohibe usar herencia', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000156', '10000000-0000-0000-0000-000000000043', 'Obliga a usar singletons', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000157', '10000000-0000-0000-0000-000000000044', 'Ocultar datos internos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000158', '10000000-0000-0000-0000-000000000044', 'Tratar objetos distintos mediante una interfaz comun con comportamientos especificos', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000159', '10000000-0000-0000-0000-000000000044', 'Duplicar codigo en subclases', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000160', '10000000-0000-0000-0000-000000000044', 'Convertir objetos a JSON', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000161', '10000000-0000-0000-0000-000000000045', 'O(n)', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000162', '10000000-0000-0000-0000-000000000045', 'O(1) promedio', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000163', '10000000-0000-0000-0000-000000000045', 'O(log n)', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000164', '10000000-0000-0000-0000-000000000045', 'O(n log n)', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000165', '10000000-0000-0000-0000-000000000046', 'Son sinonimos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000166', '10000000-0000-0000-0000-000000000046', 'Stub retorna datos fijos; mock verifica interacciones y expectativas', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000167', '10000000-0000-0000-0000-000000000046', 'Mock solo funciona en integracion', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000168', '10000000-0000-0000-0000-000000000046', 'Stub requiere Docker', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000169', '10000000-0000-0000-0000-000000000047', 'Solo validando en el frontend', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000170', '10000000-0000-0000-0000-000000000047', 'Escapando salida, sanitizando HTML y usando Content-Security-Policy', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000171', '10000000-0000-0000-0000-000000000047', 'Deshabilitando JavaScript en el navegador', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000172', '10000000-0000-0000-0000-000000000047', 'Usando solo HTTP sin TLS', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000173', '10000000-0000-0000-0000-000000000048', 'Son iguales', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000174', '10000000-0000-0000-0000-000000000048', 'Lineal O(n) recorre todo; binaria O(log n) requiere array ordenado', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000175', '10000000-0000-0000-0000-000000000048', 'Binaria solo funciona en listas enlazadas', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000176', '10000000-0000-0000-0000-000000000048', 'Lineal requiere array ordenado', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000177', '10000000-0000-0000-0000-000000000049', 'Solo como base de datos relacional', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000178', '10000000-0000-0000-0000-000000000049', 'Cache en memoria, sesiones y colas de alta velocidad', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000179', '10000000-0000-0000-0000-000000000049', 'Compilador de Go', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000180', '10000000-0000-0000-0000-000000000049', 'Balanceador de carga', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000181', '10000000-0000-0000-0000-000000000050', 'Un archivo CSS', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000182', '10000000-0000-0000-0000-000000000050', 'Unidad de UI con template, logica y estilos gestionada por Angular', true, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000183', '10000000-0000-0000-0000-000000000050', 'Una tabla de base de datos', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO question_options (id, question_id, content, is_correct, created_at, updated_at) VALUES ('20000000-0000-0000-0000-000000000184', '10000000-0000-0000-0000-000000000050', 'Un hook de React', false, NOW(), NOW()) ON CONFLICT (id) DO NOTHING;

-- Code challenges (para code_completion)
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', 'package main

func Suma(a, b int) int {
    // completa aqui
}
', '5', 'go', '[]', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000009', 'function primero<T>(arr: T[]): T | undefined {
  // completa aqui
}
', '1', 'typescript', '[]', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000013', 'function pares(arr) {
  // retorna solo los pares
}
', '[2,4]', 'javascript', '[]', NOW(), NOW()) ON CONFLICT (id) DO NOTHING;
INSERT INTO code_challenges (id, question_id, starter_code, expected_output, language, test_cases_json, created_at, updated_at) VALUES ('30000000-0000-0000-0000-000000000017', '10000000-0000-0000-0000-000000000017', 'def factorial(n):
    # completa aqui
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
