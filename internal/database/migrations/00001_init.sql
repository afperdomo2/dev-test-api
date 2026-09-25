-- +goose Up
--
-- PostgreSQL database dump
--


-- Dumped from database version 16.15
-- Dumped by pg_dump version 16.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: question_difficulty; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.question_difficulty AS ENUM (
    'beginner',
    'intermediate',
    'advanced'
);


--
-- Name: TYPE question_difficulty; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.question_difficulty IS 'Dificultad de la pregunta: beginner, intermediate o advanced';


--
-- Name: question_source; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.question_source AS ENUM (
    'ai_generated',
    'manual',
    'imported'
);


--
-- Name: TYPE question_source; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.question_source IS 'Origen de la pregunta: ai_generated, manual o imported';


--
-- Name: question_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.question_type AS ENUM (
    'single_choice',
    'multiple_choice',
    'code_completion',
    'true_false'
);


--
-- Name: TYPE question_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.question_type IS 'Tipo de pregunta: single_choice, multiple_choice, code_completion o true_false';


--
-- Name: report_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.report_status AS ENUM (
    'pending',
    'reviewed'
);


--
-- Name: TYPE report_status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.report_status IS 'Estado de un reporte: pending o reviewed';


--
-- Name: session_difficulty; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.session_difficulty AS ENUM (
    'beginner',
    'intermediate',
    'advanced'
);


--
-- Name: TYPE session_difficulty; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.session_difficulty IS 'Dificultad de una sesión: beginner, intermediate o advanced';


--
-- Name: session_mode; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.session_mode AS ENUM (
    'generate',
    'review',
    'daily'
);


--
-- Name: TYPE session_mode; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.session_mode IS 'Modo de sesión: generate (IA), review (repaso manual) o daily (sesión diaria inteligente)';


--
-- Name: session_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.session_status AS ENUM (
    'in_progress',
    'completed',
    'cancelled'
);


--
-- Name: TYPE session_status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TYPE public.session_status IS 'Estado de una sesión: in_progress, completed o cancelled';


--
-- Name: calc_session_score(); Type: FUNCTION; Schema: public; Owner: -
--

-- +goose StatementBegin
CREATE FUNCTION public.calc_session_score() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
		BEGIN
			IF NEW.status = 'completed' AND (OLD.status IS DISTINCT FROM 'completed') THEN
				NEW.score = (
					SELECT ROUND(
						(COUNT(*) FILTER (WHERE is_correct = true)::decimal /
						 GREATEST(COUNT(*), 1)) * 100, 1
					)
					FROM session_answers
					WHERE session_id = NEW.id
				);
			END IF;
			RETURN NEW;
		END;
		$$;
-- +goose StatementEnd


--
-- Name: FUNCTION calc_session_score(); Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON FUNCTION public.calc_session_score() IS 'Calcula el score (0-100) de la sesión al pasarla a completed, basado en el porcentaje de respuestas correctas';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: code_challenges; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.code_challenges (
    id uuid NOT NULL,
    question_id uuid NOT NULL,
    starter_code text,
    expected_output text,
    language character varying(50) NOT NULL,
    test_cases_json jsonb,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE code_challenges; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.code_challenges IS 'Desafío de código asociado a preguntas tipo code_completion.';


--
-- Name: COLUMN code_challenges.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.id IS 'Identificador UUID del desafío';


--
-- Name: COLUMN code_challenges.question_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.question_id IS 'Pregunta (única) asociada al desafío';


--
-- Name: COLUMN code_challenges.starter_code; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.starter_code IS 'Código inicial provisto al usuario';


--
-- Name: COLUMN code_challenges.expected_output; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.expected_output IS 'Salida esperada para validar la solución';


--
-- Name: COLUMN code_challenges.language; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.language IS 'Lenguaje de programación del desafío';


--
-- Name: COLUMN code_challenges.test_cases_json; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.test_cases_json IS 'Casos de prueba en formato JSONB';


--
-- Name: COLUMN code_challenges.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.created_at IS 'Fecha de creación del desafío';


--
-- Name: COLUMN code_challenges.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.code_challenges.updated_at IS 'Fecha de última actualización del desafío';


--
-- Name: question_options; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.question_options (
    id uuid NOT NULL,
    question_id uuid NOT NULL,
    content text NOT NULL,
    is_correct boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE question_options; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.question_options IS 'Opciones de respuesta de preguntas de tipo choice (single_choice, multiple_choice, true_false).';


--
-- Name: COLUMN question_options.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_options.id IS 'Identificador UUID de la opción';


--
-- Name: COLUMN question_options.question_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_options.question_id IS 'Pregunta a la que pertenece la opción';


--
-- Name: COLUMN question_options.content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_options.content IS 'Texto de la opción';


--
-- Name: COLUMN question_options.is_correct; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_options.is_correct IS 'Indica si la opción es la correcta';


--
-- Name: COLUMN question_options.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_options.created_at IS 'Fecha de creación de la opción';


--
-- Name: COLUMN question_options.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_options.updated_at IS 'Fecha de última actualización de la opción';


--
-- Name: question_topics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.question_topics (
    question_id uuid NOT NULL,
    topic_id uuid NOT NULL,
    created_at timestamp with time zone
);


--
-- Name: TABLE question_topics; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.question_topics IS 'Relación N:M entre preguntas y temas.';


--
-- Name: COLUMN question_topics.question_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_topics.question_id IS 'Pregunta (parte de la PK compuesta)';


--
-- Name: COLUMN question_topics.topic_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_topics.topic_id IS 'Tema (parte de la PK compuesta)';


--
-- Name: COLUMN question_topics.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.question_topics.created_at IS 'Fecha de creación de la relación';


--
-- Name: questions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.questions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    type character varying(30) NOT NULL,
    content text NOT NULL,
    explanation text,
    difficulty character varying(20) DEFAULT 'intermediate'::character varying NOT NULL,
    is_public boolean DEFAULT false NOT NULL,
    source character varying(20) DEFAULT 'ai_generated'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: TABLE questions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.questions IS 'Banco de preguntas. Las de source ai_generated y las públicas (is_public=true) son visibles para todos; el resto es privado de su creador.';


--
-- Name: COLUMN questions.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.id IS 'Identificador UUID de la pregunta';


--
-- Name: COLUMN questions.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.user_id IS 'Usuario creador de la pregunta';


--
-- Name: COLUMN questions.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.type IS 'Tipo de pregunta (ver enum question_type)';


--
-- Name: COLUMN questions.content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.content IS 'Enunciado de la pregunta';


--
-- Name: COLUMN questions.explanation; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.explanation IS 'Explicación de la respuesta correcta';


--
-- Name: COLUMN questions.difficulty; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.difficulty IS 'Dificultad: beginner, intermediate o advanced';


--
-- Name: COLUMN questions.is_public; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.is_public IS 'Indica si la pregunta es pública (visible para todos)';


--
-- Name: COLUMN questions.source; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.source IS 'Origen de la pregunta: ai_generated, manual o imported';


--
-- Name: COLUMN questions.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.created_at IS 'Fecha de creación de la pregunta';


--
-- Name: COLUMN questions.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.updated_at IS 'Fecha de última actualización de la pregunta';


--
-- Name: COLUMN questions.deleted_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.questions.deleted_at IS 'Soft delete (borrado lógico) de la pregunta';


--
-- Name: reports; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reports (
    id uuid NOT NULL,
    reporter_user_id uuid NOT NULL,
    target_type character varying(20) NOT NULL,
    target_id uuid NOT NULL,
    reason text NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE reports; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.reports IS 'Reportes de contenido (preguntas o tips) hechos por usuarios para revisión del administrador.';


--
-- Name: COLUMN reports.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.id IS 'Identificador UUID del reporte';


--
-- Name: COLUMN reports.reporter_user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.reporter_user_id IS 'Usuario que realizó el reporte';


--
-- Name: COLUMN reports.target_type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.target_type IS 'Tipo de contenido reportado: question o tip';


--
-- Name: COLUMN reports.target_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.target_id IS 'Identificador del contenido reportado';


--
-- Name: COLUMN reports.reason; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.reason IS 'Motivo del reporte en texto libre';


--
-- Name: COLUMN reports.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.status IS 'Estado del reporte: pending o reviewed';


--
-- Name: COLUMN reports.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.created_at IS 'Fecha de creación del reporte';


--
-- Name: COLUMN reports.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.reports.updated_at IS 'Fecha de última actualización del reporte';


--
-- Name: session_answers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.session_answers (
    id uuid NOT NULL,
    session_id uuid NOT NULL,
    user_id uuid NOT NULL,
    question_id uuid NOT NULL,
    answer_text text,
    selected_options jsonb,
    is_correct boolean DEFAULT false NOT NULL,
    ai_feedback text,
    response_time_ms bigint,
    created_at timestamp with time zone
);


--
-- Name: TABLE session_answers; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.session_answers IS 'Respuestas del usuario a preguntas dentro de una sesión (con resultado y tiempo de respuesta).';


--
-- Name: COLUMN session_answers.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.id IS 'Identificador UUID de la respuesta';


--
-- Name: COLUMN session_answers.session_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.session_id IS 'Sesión a la que pertenece la respuesta';


--
-- Name: COLUMN session_answers.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.user_id IS 'Usuario que respondió';


--
-- Name: COLUMN session_answers.question_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.question_id IS 'Pregunta respondida';


--
-- Name: COLUMN session_answers.answer_text; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.answer_text IS 'Respuesta en texto (para code_completion)';


--
-- Name: COLUMN session_answers.selected_options; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.selected_options IS 'Opciones seleccionadas en formato JSONB (array de UUIDs)';


--
-- Name: COLUMN session_answers.is_correct; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.is_correct IS 'Indica si la respuesta fue correcta';


--
-- Name: COLUMN session_answers.ai_feedback; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.ai_feedback IS 'Feedback generado por IA sobre la respuesta';


--
-- Name: COLUMN session_answers.response_time_ms; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.response_time_ms IS 'Tiempo de respuesta en milisegundos';


--
-- Name: COLUMN session_answers.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_answers.created_at IS 'Fecha de creación de la respuesta';


--
-- Name: session_topics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.session_topics (
    session_id uuid NOT NULL,
    topic_id uuid NOT NULL,
    created_at timestamp with time zone
);


--
-- Name: TABLE session_topics; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.session_topics IS 'Relación N:M entre sesiones y temas.';


--
-- Name: COLUMN session_topics.session_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_topics.session_id IS 'Sesión (parte de la PK compuesta)';


--
-- Name: COLUMN session_topics.topic_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_topics.topic_id IS 'Tema (parte de la PK compuesta)';


--
-- Name: COLUMN session_topics.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.session_topics.created_at IS 'Fecha de creación de la relación';


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    name character varying(200) NOT NULL,
    status character varying(20) DEFAULT 'in_progress'::character varying NOT NULL,
    mode character varying(20) DEFAULT 'generate'::character varying NOT NULL,
    difficulty character varying(20) DEFAULT 'intermediate'::character varying NOT NULL,
    question_limit bigint,
    questions_generated bigint DEFAULT 0 NOT NULL,
    tip_id uuid,
    score numeric,
    started_at timestamp with time zone,
    finished_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE sessions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.sessions IS 'Sesiones de estudio. Modos: generate (IA), review (repaso manual) y daily (sesión diaria inteligente con tip).';


--
-- Name: COLUMN sessions.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.id IS 'Identificador UUID de la sesión';


--
-- Name: COLUMN sessions.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.user_id IS 'Usuario al que pertenece la sesión';


--
-- Name: COLUMN sessions.name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.name IS 'Nombre visible de la sesión';


--
-- Name: COLUMN sessions.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.status IS 'Estado: in_progress, completed o cancelled';


--
-- Name: COLUMN sessions.mode; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.mode IS 'Modo de sesión (ver enum session_mode)';


--
-- Name: COLUMN sessions.difficulty; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.difficulty IS 'Dificultad de las preguntas de la sesión (ver enum session_difficulty)';


--
-- Name: COLUMN sessions.question_limit; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.question_limit IS 'Límite de preguntas de la sesión (NULL = sin límite)';


--
-- Name: COLUMN sessions.questions_generated; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.questions_generated IS 'Cantidad de preguntas generadas para la sesión';


--
-- Name: COLUMN sessions.tip_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.tip_id IS 'Tip del día asociado a una sesión en modo daily (nullable en el resto de modos)';


--
-- Name: COLUMN sessions.score; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.score IS 'Puntuación final calculada al completar la sesión (0-100, ver trigger trg_sessions_score)';


--
-- Name: COLUMN sessions.started_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.started_at IS 'Fecha de inicio de la sesión';


--
-- Name: COLUMN sessions.finished_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.finished_at IS 'Fecha de finalización de la sesión (NULL mientras esté en progreso)';


--
-- Name: COLUMN sessions.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.created_at IS 'Fecha de creación de la sesión';


--
-- Name: COLUMN sessions.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.sessions.updated_at IS 'Fecha de última actualización de la sesión';


--
-- Name: tips; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tips (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    topic_id uuid NOT NULL,
    content text NOT NULL,
    source character varying(20) DEFAULT 'manual'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: TABLE tips; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.tips IS 'Tarjetas tip creadas por administradores (siempre públicas). Leídas por cualquier usuario autenticado y usadas por la sesión diaria.';


--
-- Name: COLUMN tips.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.id IS 'Identificador UUID del tip';


--
-- Name: COLUMN tips.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.user_id IS 'Administrador que creó el tip';


--
-- Name: COLUMN tips.topic_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.topic_id IS 'Tema al que pertenece el tip (un único tema por tip)';


--
-- Name: COLUMN tips.content; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.content IS 'Contenido del tip en texto plano';


--
-- Name: COLUMN tips.source; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.source IS 'Origen del tip: manual o imported (sin contenido generado por IA)';


--
-- Name: COLUMN tips.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.created_at IS 'Fecha de creación del tip';


--
-- Name: COLUMN tips.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.updated_at IS 'Fecha de última actualización del tip';


--
-- Name: COLUMN tips.deleted_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.tips.deleted_at IS 'Soft delete (borrado lógico) del tip';


--
-- Name: topics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.topics (
    id uuid NOT NULL,
    slug character varying(100) NOT NULL,
    name character varying(150) NOT NULL,
    category character varying(50) NOT NULL,
    is_system boolean DEFAULT false NOT NULL,
    created_by uuid,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE topics; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.topics IS 'Temas de estudio. Los temas del sistema tienen is_system=true; los creados por usuarios quedan ligados a created_by.';


--
-- Name: COLUMN topics.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.id IS 'Identificador UUID del tema';


--
-- Name: COLUMN topics.slug; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.slug IS 'Identificador corto único por creador (system o usuario)';


--
-- Name: COLUMN topics.name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.name IS 'Nombre visible del tema';


--
-- Name: COLUMN topics.category; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.category IS 'Categoría a la que pertenece el tema';


--
-- Name: COLUMN topics.is_system; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.is_system IS 'Indica si el tema es del sistema (visible para todos)';


--
-- Name: COLUMN topics.created_by; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.created_by IS 'Usuario creador del tema (NULL si es del sistema)';


--
-- Name: COLUMN topics.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.created_at IS 'Fecha de creación del tema';


--
-- Name: COLUMN topics.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.topics.updated_at IS 'Fecha de última actualización del tema';


--
-- Name: user_question_progress; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_question_progress (
    user_id uuid NOT NULL,
    question_id uuid NOT NULL,
    repetitions bigint DEFAULT 0 NOT NULL,
    ease_factor numeric DEFAULT 2.5 NOT NULL,
    interval_days bigint DEFAULT 0 NOT NULL,
    next_review_at timestamp with time zone,
    last_reviewed_at timestamp with time zone,
    is_saved boolean DEFAULT false NOT NULL,
    is_mastered boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE user_question_progress; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.user_question_progress IS 'Progreso individual del usuario por pregunta (repeticiones, facilidad, fecha de próximo repaso) basado en SM-2.';


--
-- Name: COLUMN user_question_progress.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.user_id IS 'Usuario (parte de la PK compuesta)';


--
-- Name: COLUMN user_question_progress.question_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.question_id IS 'Pregunta (parte de la PK compuesta)';


--
-- Name: COLUMN user_question_progress.repetitions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.repetitions IS 'Número de repeticiones exitosas consecutivas';


--
-- Name: COLUMN user_question_progress.ease_factor; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.ease_factor IS 'Factor de facilidad SM-2 (default 2.5)';


--
-- Name: COLUMN user_question_progress.interval_days; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.interval_days IS 'Intervalo actual entre repasos en días';


--
-- Name: COLUMN user_question_progress.next_review_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.next_review_at IS 'Fecha del próximo repaso programado';


--
-- Name: COLUMN user_question_progress.last_reviewed_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.last_reviewed_at IS 'Fecha del último repaso realizado';


--
-- Name: COLUMN user_question_progress.is_saved; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.is_saved IS 'Indica si el usuario guardó la pregunta como favorita';


--
-- Name: COLUMN user_question_progress.is_mastered; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.is_mastered IS 'Indica si la pregunta se considera dominada (no vuelve a aparecer en repasos)';


--
-- Name: COLUMN user_question_progress.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.created_at IS 'Fecha de creación del progreso';


--
-- Name: COLUMN user_question_progress.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_question_progress.updated_at IS 'Fecha de última actualización del progreso';


--
-- Name: user_streaks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_streaks (
    user_id uuid NOT NULL,
    current_streak bigint DEFAULT 0 NOT NULL,
    longest_streak bigint DEFAULT 0 NOT NULL,
    last_activity_at timestamp with time zone,
    last_activity_date character varying(10),
    freezes_available bigint DEFAULT 0 NOT NULL,
    daily_xp_today bigint DEFAULT 0 NOT NULL,
    daily_xp_date character varying(10),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE user_streaks; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.user_streaks IS 'Racha diaria de actividad del usuario (días consecutivos con pregunta respondida), racha más larga y freezes disponibles.';


--
-- Name: COLUMN user_streaks.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.user_id IS 'Usuario (PK)';


--
-- Name: COLUMN user_streaks.current_streak; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.current_streak IS 'Racha actual en días consecutivos';


--
-- Name: COLUMN user_streaks.longest_streak; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.longest_streak IS 'Racha más larga alcanzada en días';


--
-- Name: COLUMN user_streaks.last_activity_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.last_activity_at IS 'Timestamp de la última actividad';


--
-- Name: COLUMN user_streaks.last_activity_date; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.last_activity_date IS 'Fecha (YYYY-MM-DD, UTC) de la última actividad para comparar días consecutivos';


--
-- Name: COLUMN user_streaks.freezes_available; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.freezes_available IS 'Freezes disponibles (congelan un día sin actividad, tope 3)';


--
-- Name: COLUMN user_streaks.daily_xp_today; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.daily_xp_today IS 'XP acumulado en el día actual para la meta diaria';


--
-- Name: COLUMN user_streaks.daily_xp_date; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.daily_xp_date IS 'Fecha (YYYY-MM-DD, UTC) a la que corresponde daily_xp_today';


--
-- Name: COLUMN user_streaks.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.created_at IS 'Fecha de creación del registro';


--
-- Name: COLUMN user_streaks.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_streaks.updated_at IS 'Fecha de última actualización del registro';


--
-- Name: user_topic_xps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_topic_xps (
    user_id uuid NOT NULL,
    topic_id uuid NOT NULL,
    xp bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE user_topic_xps; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.user_topic_xps IS 'XP del usuario por tema (desglose del XP global).';


--
-- Name: COLUMN user_topic_xps.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_topic_xps.user_id IS 'Usuario (parte de la PK compuesta)';


--
-- Name: COLUMN user_topic_xps.topic_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_topic_xps.topic_id IS 'Tema (parte de la PK compuesta)';


--
-- Name: COLUMN user_topic_xps.xp; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_topic_xps.xp IS 'XP acumulado en el tema';


--
-- Name: COLUMN user_topic_xps.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_topic_xps.created_at IS 'Fecha de creación del registro';


--
-- Name: COLUMN user_topic_xps.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_topic_xps.updated_at IS 'Fecha de última actualización del registro';


--
-- Name: user_xps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_xps (
    user_id uuid NOT NULL,
    xp bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


--
-- Name: TABLE user_xps; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.user_xps IS 'XP global del usuario (suma acumulada de todas las respuestas correctas).';


--
-- Name: COLUMN user_xps.user_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_xps.user_id IS 'Usuario (PK)';


--
-- Name: COLUMN user_xps.xp; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_xps.xp IS 'XP acumulado total';


--
-- Name: COLUMN user_xps.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_xps.created_at IS 'Fecha de creación del registro';


--
-- Name: COLUMN user_xps.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_xps.updated_at IS 'Fecha de última actualización del registro';


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    is_admin boolean DEFAULT false,
    daily_import_limit bigint DEFAULT 200 NOT NULL,
    daily_ai_limit bigint DEFAULT 20 NOT NULL,
    daily_goal_xp bigint DEFAULT 50 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


--
-- Name: TABLE users; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.users IS 'Usuarios de la plataforma (login, permisos y límites diarios)';


--
-- Name: COLUMN users.id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.id IS 'Identificador UUID del usuario';


--
-- Name: COLUMN users.email; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.email IS 'Email único del usuario (usado para autenticación)';


--
-- Name: COLUMN users.password_hash; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.password_hash IS 'Hash del password (BCrypt); nunca se expone en JSON';


--
-- Name: COLUMN users.is_admin; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.is_admin IS 'Indica si el usuario tiene rol administrador';


--
-- Name: COLUMN users.daily_import_limit; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.daily_import_limit IS 'Cupo diario de preguntas importadas por CSV';


--
-- Name: COLUMN users.daily_ai_limit; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.daily_ai_limit IS 'Cupo diario de preguntas generadas por IA';


--
-- Name: COLUMN users.daily_goal_xp; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.daily_goal_xp IS 'Meta diaria de XP configurable por el usuario (presets: casual=25, regular=50, serious=100)';


--
-- Name: COLUMN users.created_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.created_at IS 'Fecha de creación del usuario';


--
-- Name: COLUMN users.updated_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.updated_at IS 'Fecha de última actualización del usuario';


--
-- Name: COLUMN users.deleted_at; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.deleted_at IS 'Soft delete (borrado lógico) del usuario';


--
-- Name: code_challenges code_challenges_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.code_challenges
    ADD CONSTRAINT code_challenges_pkey PRIMARY KEY (id);


--
-- Name: question_options question_options_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.question_options
    ADD CONSTRAINT question_options_pkey PRIMARY KEY (id);


--
-- Name: question_topics question_topics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.question_topics
    ADD CONSTRAINT question_topics_pkey PRIMARY KEY (question_id, topic_id);


--
-- Name: questions questions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT questions_pkey PRIMARY KEY (id);


--
-- Name: reports reports_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reports
    ADD CONSTRAINT reports_pkey PRIMARY KEY (id);


--
-- Name: session_answers session_answers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session_answers
    ADD CONSTRAINT session_answers_pkey PRIMARY KEY (id);


--
-- Name: session_topics session_topics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session_topics
    ADD CONSTRAINT session_topics_pkey PRIMARY KEY (session_id, topic_id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: tips tips_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tips
    ADD CONSTRAINT tips_pkey PRIMARY KEY (id);


--
-- Name: topics topics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT topics_pkey PRIMARY KEY (id);


--
-- Name: user_question_progress user_question_progress_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_question_progress
    ADD CONSTRAINT user_question_progress_pkey PRIMARY KEY (user_id, question_id);


--
-- Name: user_streaks user_streaks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_streaks
    ADD CONSTRAINT user_streaks_pkey PRIMARY KEY (user_id);


--
-- Name: user_topic_xps user_topic_xps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_topic_xps
    ADD CONSTRAINT user_topic_xps_pkey PRIMARY KEY (user_id, topic_id);


--
-- Name: user_xps user_xps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_xps
    ADD CONSTRAINT user_xps_pkey PRIMARY KEY (user_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_code_challenges_question_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_code_challenges_question_id ON public.code_challenges USING btree (question_id);


--
-- Name: idx_question_options_question_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_question_options_question_id ON public.question_options USING btree (question_id);


--
-- Name: idx_questions_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_questions_deleted_at ON public.questions USING btree (deleted_at);


--
-- Name: idx_questions_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_questions_type ON public.questions USING btree (type);


--
-- Name: idx_questions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_questions_user_id ON public.questions USING btree (user_id);


--
-- Name: idx_reports_reporter_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reports_reporter_user_id ON public.reports USING btree (reporter_user_id);


--
-- Name: idx_reports_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reports_status ON public.reports USING btree (status);


--
-- Name: idx_reports_target; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reports_target ON public.reports USING btree (target_type, target_id);


--
-- Name: idx_reports_target_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reports_target_id ON public.reports USING btree (target_id);


--
-- Name: idx_session_answers_question_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_session_answers_question_id ON public.session_answers USING btree (question_id);


--
-- Name: idx_session_answers_session; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_session_answers_session ON public.session_answers USING btree (session_id);


--
-- Name: idx_session_answers_session_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_session_answers_session_id ON public.session_answers USING btree (session_id);


--
-- Name: idx_session_answers_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_session_answers_user_id ON public.session_answers USING btree (user_id);


--
-- Name: idx_sessions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_user_id ON public.sessions USING btree (user_id);


--
-- Name: idx_sessions_user_started; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_user_started ON public.sessions USING btree (user_id, started_at DESC);


--
-- Name: idx_sessions_user_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_user_status ON public.sessions USING btree (user_id, status);


--
-- Name: idx_tips_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tips_deleted_at ON public.tips USING btree (deleted_at);


--
-- Name: idx_tips_topic_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tips_topic_id ON public.tips USING btree (topic_id);


--
-- Name: idx_tips_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tips_user_id ON public.tips USING btree (user_id);


--
-- Name: idx_topics_slug_created_by; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_topics_slug_created_by ON public.topics USING btree (slug, created_by) NULLS NOT DISTINCT;


--
-- Name: idx_user_question_progress_mastered; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_question_progress_mastered ON public.user_question_progress USING btree (user_id, is_mastered) WHERE (is_mastered = false);


--
-- Name: idx_user_question_progress_next_review; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_question_progress_next_review ON public.user_question_progress USING btree (user_id, next_review_at) WHERE (next_review_at IS NOT NULL);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: sessions trg_sessions_score; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_sessions_score BEFORE UPDATE ON public.sessions FOR EACH ROW EXECUTE FUNCTION public.calc_session_score();


--
-- Name: TRIGGER trg_sessions_score ON sessions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TRIGGER trg_sessions_score ON public.sessions IS 'Ejecuta calc_session_score() al actualizar una sesión para asignar el score al completarla';


--
-- Name: code_challenges fk_questions_code_challenge; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.code_challenges
    ADD CONSTRAINT fk_questions_code_challenge FOREIGN KEY (question_id) REFERENCES public.questions(id);


--
-- Name: question_options fk_questions_options; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.question_options
    ADD CONSTRAINT fk_questions_options FOREIGN KEY (question_id) REFERENCES public.questions(id);


--
-- Name: question_topics fk_questions_question_topics; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.question_topics
    ADD CONSTRAINT fk_questions_question_topics FOREIGN KEY (question_id) REFERENCES public.questions(id);


--
-- Name: questions fk_questions_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.questions
    ADD CONSTRAINT fk_questions_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: reports fk_reports_reporter; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reports
    ADD CONSTRAINT fk_reports_reporter FOREIGN KEY (reporter_user_id) REFERENCES public.users(id);


--
-- Name: session_answers fk_session_answers_question; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session_answers
    ADD CONSTRAINT fk_session_answers_question FOREIGN KEY (question_id) REFERENCES public.questions(id);


--
-- Name: session_answers fk_session_answers_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session_answers
    ADD CONSTRAINT fk_session_answers_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: session_answers fk_sessions_answers; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session_answers
    ADD CONSTRAINT fk_sessions_answers FOREIGN KEY (session_id) REFERENCES public.sessions(id);


--
-- Name: session_topics fk_sessions_session_topics; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.session_topics
    ADD CONSTRAINT fk_sessions_session_topics FOREIGN KEY (session_id) REFERENCES public.sessions(id);


--
-- Name: sessions fk_sessions_tip; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT fk_sessions_tip FOREIGN KEY (tip_id) REFERENCES public.tips(id);


--
-- Name: sessions fk_sessions_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: tips fk_tips_topic; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tips
    ADD CONSTRAINT fk_tips_topic FOREIGN KEY (topic_id) REFERENCES public.topics(id);


--
-- Name: tips fk_tips_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tips
    ADD CONSTRAINT fk_tips_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: topics fk_topics_creator; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT fk_topics_creator FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: user_question_progress fk_user_question_progress_question; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_question_progress
    ADD CONSTRAINT fk_user_question_progress_question FOREIGN KEY (question_id) REFERENCES public.questions(id);


--
-- Name: user_question_progress fk_user_question_progress_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_question_progress
    ADD CONSTRAINT fk_user_question_progress_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_streaks fk_user_streaks_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_streaks
    ADD CONSTRAINT fk_user_streaks_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_topic_xps fk_user_topic_xps_topic; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_topic_xps
    ADD CONSTRAINT fk_user_topic_xps_topic FOREIGN KEY (topic_id) REFERENCES public.topics(id);


--
-- Name: user_topic_xps fk_user_topic_xps_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_topic_xps
    ADD CONSTRAINT fk_user_topic_xps_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_xps fk_user_xps_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_xps
    ADD CONSTRAINT fk_user_xps_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--



-- +goose Down
DROP TRIGGER IF EXISTS trg_sessions_score ON public.sessions;
DROP FUNCTION IF EXISTS public.calc_session_score();
DROP TABLE IF EXISTS public.user_question_progress;
DROP TABLE IF EXISTS public.user_topic_xps;
DROP TABLE IF EXISTS public.user_xps;
DROP TABLE IF EXISTS public.user_streaks;
DROP TABLE IF EXISTS public.reports;
DROP TABLE IF EXISTS public.session_answers;
DROP TABLE IF EXISTS public.session_topics;
DROP TABLE IF EXISTS public.sessions;
DROP TABLE IF EXISTS public.tips;
DROP TABLE IF EXISTS public.code_challenges;
DROP TABLE IF EXISTS public.question_options;
DROP TABLE IF EXISTS public.question_topics;
DROP TABLE IF EXISTS public.questions;
DROP TABLE IF EXISTS public.topics;
DROP TABLE IF EXISTS public.users;
DROP TYPE IF EXISTS public.question_difficulty;
DROP TYPE IF EXISTS public.question_source;
DROP TYPE IF EXISTS public.question_type;
DROP TYPE IF EXISTS public.report_status;
DROP TYPE IF EXISTS public.session_difficulty;
DROP TYPE IF EXISTS public.session_mode;
DROP TYPE IF EXISTS public.session_status;
