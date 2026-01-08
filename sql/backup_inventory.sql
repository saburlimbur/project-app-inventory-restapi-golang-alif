--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: TABLE categories; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.categories IS 'Tabel untuk kategori barang';


--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
    id integer NOT NULL,
    category_id integer NOT NULL,
    rack_id integer,
    sku character varying(50) NOT NULL,
    name character varying(200) NOT NULL,
    description text,
    unit character varying(20) DEFAULT 'pcs'::character varying NOT NULL,
    price numeric(15,2) NOT NULL,
    cost numeric(15,2) DEFAULT 0,
    stock integer DEFAULT 0 NOT NULL,
    minimum_stock integer DEFAULT 5 NOT NULL,
    weight numeric(10,2) DEFAULT 0,
    dimensions character varying(50),
    is_active boolean DEFAULT true,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT items_cost_check CHECK ((cost >= (0)::numeric)),
    CONSTRAINT items_price_check CHECK ((price >= (0)::numeric)),
    CONSTRAINT items_stock_check CHECK ((stock >= 0))
);


ALTER TABLE public.items OWNER TO postgres;

--
-- Name: TABLE items; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.items IS 'Tabel untuk barang/produk';


--
-- Name: items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.items_id_seq OWNER TO postgres;

--
-- Name: items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.items_id_seq OWNED BY public.items.id;


--
-- Name: racks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.racks (
    id integer NOT NULL,
    warehouse_id integer NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(100) NOT NULL,
    location character varying(100),
    capacity integer DEFAULT 0,
    description text,
    is_active boolean DEFAULT true,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.racks OWNER TO postgres;

--
-- Name: TABLE racks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.racks IS 'Tabel untuk rak penyimpanan di gudang';


--
-- Name: racks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.racks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.racks_id_seq OWNER TO postgres;

--
-- Name: racks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.racks_id_seq OWNED BY public.racks.id;


--
-- Name: sale_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sale_items (
    id integer NOT NULL,
    sale_id integer NOT NULL,
    item_id integer NOT NULL,
    quantity integer NOT NULL,
    unit_price numeric(15,2) NOT NULL,
    subtotal numeric(15,2) NOT NULL,
    discount numeric(15,2) DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sale_items_discount_check CHECK ((discount >= (0)::numeric)),
    CONSTRAINT sale_items_quantity_check CHECK ((quantity > 0)),
    CONSTRAINT sale_items_subtotal_check CHECK ((subtotal >= (0)::numeric)),
    CONSTRAINT sale_items_unit_price_check CHECK ((unit_price >= (0)::numeric))
);


ALTER TABLE public.sale_items OWNER TO postgres;

--
-- Name: TABLE sale_items; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sale_items IS 'Tabel untuk detail item penjualan';


--
-- Name: sale_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sale_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sale_items_id_seq OWNER TO postgres;

--
-- Name: sale_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sale_items_id_seq OWNED BY public.sale_items.id;


--
-- Name: sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sales (
    id integer NOT NULL,
    invoice_number character varying(50) NOT NULL,
    customer_name character varying(100),
    customer_phone character varying(20),
    customer_email character varying(100),
    sale_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    total_amount numeric(15,2) DEFAULT 0 NOT NULL,
    discount numeric(15,2) DEFAULT 0,
    tax numeric(15,2) DEFAULT 0,
    grand_total numeric(15,2) DEFAULT 0 NOT NULL,
    payment_method character varying(50),
    payment_status character varying(20) DEFAULT 'pending'::character varying,
    notes text,
    created_by integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sales_discount_check CHECK ((discount >= (0)::numeric)),
    CONSTRAINT sales_grand_total_check CHECK ((grand_total >= (0)::numeric)),
    CONSTRAINT sales_payment_status_check CHECK (((payment_status)::text = ANY ((ARRAY['pending'::character varying, 'paid'::character varying, 'cancelled'::character varying])::text[]))),
    CONSTRAINT sales_tax_check CHECK ((tax >= (0)::numeric)),
    CONSTRAINT sales_total_amount_check CHECK ((total_amount >= (0)::numeric))
);


ALTER TABLE public.sales OWNER TO postgres;

--
-- Name: TABLE sales; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sales IS 'Tabel untuk transaksi penjualan';


--
-- Name: sales_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sales_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sales_id_seq OWNER TO postgres;

--
-- Name: sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sales_id_seq OWNED BY public.sales.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    token uuid DEFAULT gen_random_uuid() NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    expired_at timestamp without time zone NOT NULL,
    revoked_at timestamp without time zone,
    last_activity timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    ip_address character varying(50),
    user_agent text
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: TABLE sessions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sessions IS 'Tabel untuk menyimpan session token pengguna';


--
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sessions_id_seq OWNER TO postgres;

--
-- Name: sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sessions_id_seq OWNED BY public.sessions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash character varying(255) NOT NULL,
    full_name character varying(100) NOT NULL,
    role character varying(20) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['super_admin'::character varying, 'admin'::character varying, 'staff'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: TABLE users; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.users IS 'Tabel untuk menyimpan data pengguna sistem';


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: warehouses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.warehouses (
    id integer NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(100) NOT NULL,
    address text,
    city character varying(50),
    province character varying(50),
    postal_code character varying(10),
    phone character varying(20),
    is_active boolean DEFAULT true,
    created_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.warehouses OWNER TO postgres;

--
-- Name: TABLE warehouses; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.warehouses IS 'Tabel untuk gudang penyimpanan';


--
-- Name: v_low_stock_items; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.v_low_stock_items AS
 SELECT i.id,
    i.sku,
    i.name,
    c.name AS category_name,
    r.name AS rack_name,
    w.name AS warehouse_name,
    i.stock,
    i.minimum_stock,
    i.price,
    (i.minimum_stock - i.stock) AS stock_shortage
   FROM (((public.items i
     LEFT JOIN public.categories c ON ((i.category_id = c.id)))
     LEFT JOIN public.racks r ON ((i.rack_id = r.id)))
     LEFT JOIN public.warehouses w ON ((r.warehouse_id = w.id)))
  WHERE ((i.stock < i.minimum_stock) AND (i.is_active = true))
  ORDER BY (i.minimum_stock - i.stock) DESC;


ALTER VIEW public.v_low_stock_items OWNER TO postgres;

--
-- Name: v_sales_report; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.v_sales_report AS
 SELECT s.id,
    s.invoice_number,
    s.sale_date,
    s.customer_name,
    s.grand_total,
    s.payment_status,
    u.full_name AS created_by_name,
    count(si.id) AS total_items
   FROM ((public.sales s
     LEFT JOIN public.sale_items si ON ((s.id = si.sale_id)))
     LEFT JOIN public.users u ON ((s.created_by = u.id)))
  GROUP BY s.id, s.invoice_number, s.sale_date, s.customer_name, s.grand_total, s.payment_status, u.full_name;


ALTER VIEW public.v_sales_report OWNER TO postgres;

--
-- Name: warehouses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.warehouses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.warehouses_id_seq OWNER TO postgres;

--
-- Name: warehouses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.warehouses_id_seq OWNED BY public.warehouses.id;


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items ALTER COLUMN id SET DEFAULT nextval('public.items_id_seq'::regclass);


--
-- Name: racks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks ALTER COLUMN id SET DEFAULT nextval('public.racks_id_seq'::regclass);


--
-- Name: sale_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items ALTER COLUMN id SET DEFAULT nextval('public.sale_items_id_seq'::regclass);


--
-- Name: sales id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales ALTER COLUMN id SET DEFAULT nextval('public.sales_id_seq'::regclass);


--
-- Name: sessions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions ALTER COLUMN id SET DEFAULT nextval('public.sessions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: warehouses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses ALTER COLUMN id SET DEFAULT nextval('public.warehouses_id_seq'::regclass);


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, code, name, description, is_active, created_by, created_at, updated_at) FROM stdin;
1	ELC001	Elektronik	Ini adalah kategori elektronik	t	1	2026-01-02 17:31:41.628847	2026-01-02 17:31:41.628847
2	ATK001	Alat Tulis Kantor	Kategori untuk perlengkapan alat tulis dan kebutuhan perkantoran	t	1	2026-01-02 17:39:36.361915	2026-01-02 17:39:36.361915
3	FUR001	Furniture	Kategori untuk perabotan seperti meja, kursi, dan lemari	t	1	2026-01-02 17:43:13.088278	2026-01-02 17:43:13.088278
4	TS001	Test categorys updated	Kategori untuk testing method saja	f	5	2026-01-03 20:30:38.971875	2026-01-03 21:42:30.396664
\.


--
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.items (id, category_id, rack_id, sku, name, description, unit, price, cost, stock, minimum_stock, weight, dimensions, is_active, created_by, created_at, updated_at) FROM stdin;
4	2	3	ASL-003-ASAL-A5	Asal Tulis	Buku asal tulis ukuran A5, 40 lembar, kertas HVS	pcs	22000.00	2500.00	15	30	0.00	21x15x1	f	1	2026-01-04 20:11:38.407523	2026-01-04 20:22:43.927946
2	2	3	ATK-002-PULPEN-BIRU	Pulpen Gel Biru	Pulpen gel tinta biru, ujung 0.5 mm, nyaman untuk menulis harian	pcs	4500.00	2800.00	290	60	0.00	14x1x1	t	1	2026-01-04 19:21:19.857357	2026-01-05 01:40:57.489985
3	2	3	ATK-003-BUKU-A5	Buku Tulis A5 40 Lembar Update	Buku tulis ukuran A5, 40 lembar, kertas HVS	pcs	12000.00	8500.00	130	30	0.15	21x15x1	t	1	2026-01-04 19:21:31.602527	2026-01-08 00:59:51.560284
1	2	3	ATK-001-PULPEN	Pulpen Gel Hitam	Pulpen gel tinta hitam, ujung 0.5 mm, cocok untuk kebutuhan kantor	pcs	5000.00	3000.00	220	50	0.00	14x1x1	t	1	2026-01-04 19:05:18.222682	2026-01-08 00:59:51.560284
\.


--
-- Data for Name: racks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.racks (id, warehouse_id, code, name, location, capacity, description, is_active, created_by, created_at, updated_at) FROM stdin;
1	1	RCK-001	Rack A	Lantai 1 - Zona A	100	Rak untuk menyimpan barang elektronik	t	1	2026-01-04 01:23:27.526103	2026-01-04 01:23:27.526103
2	1	RCK-002	Rack B	Lantai 1 - Zona B	150	Rak untuk menyimpan barang furniture	t	1	2026-01-04 01:23:47.789956	2026-01-04 01:23:47.789956
3	2	RCK-003	Rack C update	Lantai 2 - Zona C	210	Rak untuk menyimpan barang alat tulis update	t	1	2026-01-04 01:24:02.445008	2026-01-04 01:27:07.606347
4	2	JKW-001	Rack Test	Lantai 2 - Zona Merah	200	Rak untuk test saja	f	1	2026-01-04 01:28:30.444096	2026-01-04 01:31:11.705522
\.


--
-- Data for Name: sale_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sale_items (id, sale_id, item_id, quantity, unit_price, subtotal, discount, created_at) FROM stdin;
3	7	1	20	5000.00	100000.00	0.00	2026-01-05 01:40:57.489985
4	7	2	10	4500.00	40000.00	5000.00	2026-01-05 01:40:57.489985
5	8	3	20	12000.00	240000.00	0.00	2026-01-08 00:59:51.560284
6	8	1	10	5000.00	45000.00	5000.00	2026-01-08 00:59:51.560284
\.


--
-- Data for Name: sales; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sales (id, invoice_number, customer_name, customer_phone, customer_email, sale_date, total_amount, discount, tax, grand_total, payment_method, payment_status, notes, created_by, created_at, updated_at) FROM stdin;
1	INV-20260104-001	Budi Santoso	081234567890	budi@example.com	2026-01-04 21:27:14.980615	1500000.00	50000.00	150000.00	1600000.00	cash	pending	Pembelian langsung di toko	1	2026-01-04 21:27:14.980615	2026-01-04 21:27:14.980615
3	INV-20260104-002	Budi Santoso	\N	\N	2026-01-04 23:53:16.973214	1500000.00	50000.00	150000.00	1600000.00	cash	pending	Pembelian langsung di toko	1	2026-01-04 23:53:16.973214	2026-01-04 23:53:16.973214
4	INV-20260105-001	Siti Aminah	089876543210	siti.aminah@mail.com	2026-01-04 23:55:49.473379	2750000.00	250000.00	275000.00	2775000.00	transfer	pending	Pembayaran via transfer BCA	1	2026-01-04 23:55:49.473379	2026-01-04 23:55:49.473379
7	INV-002	Siti	081234567890	siti@mail.com	2026-01-05 01:40:57.489985	140000.00	0.00	10000.00	150000.00	cash	pending	Pembelian offline	1	2026-01-05 01:40:57.489985	2026-01-05 01:40:57.489985
8	INV-003	Maesaroh	089999999999	maesaroh_update@mail.com	2026-01-08 00:59:51.560284	285000.00	0.00	10000.00	295000.00	transfer	paid	Customer ganti metode bayar	1	2026-01-08 00:59:51.560284	2026-01-08 01:32:21.462285
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (id, user_id, token, created_at, expired_at, revoked_at, last_activity, ip_address, user_agent) FROM stdin;
1	2	8302290e-86a0-4684-9f0f-c4177c761bcb	2025-12-30 17:15:44.254723	2025-12-31 17:15:44.237179	\N	2025-12-30 17:15:44.254723	[::1]:51247	PostmanRuntime/7.51.0
2	2	8caf2812-cb7b-4246-9f53-019953738d50	2025-12-30 17:16:50.296023	2025-12-31 17:16:50.295164	\N	2025-12-30 17:16:50.296023	[::1]:51330	PostmanRuntime/7.51.0
3	2	a8c1e4b2-2935-42b0-b607-2e556aa9cfeb	2025-12-30 17:22:59.58871	2025-12-31 17:22:59.588111	\N	2025-12-30 17:22:59.58871	[::1]:51496	PostmanRuntime/7.51.0
4	2	f2e50ce4-dc5c-4092-94e9-c71b3c2e9daf	2025-12-30 17:50:22.750076	2025-12-31 17:50:22.744398	\N	2025-12-30 17:50:22.750076	[::1]:52621	PostmanRuntime/7.51.0
5	1	1d48e928-d1bd-4a28-851b-8fbb0b4214e5	2025-12-30 18:00:07.322294	2025-12-31 18:00:07.305408	\N	2025-12-30 18:00:07.322294	[::1]:52932	PostmanRuntime/7.51.0
6	1	16b3b9c6-1e6e-4aff-a646-f40f19a7a010	2025-12-30 18:01:02.318461	2025-12-31 18:01:02.317517	\N	2025-12-30 18:01:02.318461	[::1]:52932	PostmanRuntime/7.51.0
7	2	9168d087-a15f-4305-a514-9bc501f4ed10	2025-12-31 13:58:18.335875	2026-01-01 13:58:18.3339	\N	2025-12-31 13:58:18.335875	[::1]:57108	PostmanRuntime/7.51.0
8	1	4caf5106-1bfa-43fd-b943-abe234055206	2025-12-31 13:58:28.264345	2026-01-01 13:58:28.263674	\N	2025-12-31 13:58:28.264345	[::1]:57108	PostmanRuntime/7.51.0
9	2	203f782f-270a-45fd-9368-e1ab8b38fd20	2026-01-01 15:32:23.090024	2026-01-02 15:32:23.083586	\N	2026-01-01 15:32:23.090024	[::1]:51677	PostmanRuntime/7.51.0
10	1	44efe7d1-0292-4778-9396-4ceef858c95f	2026-01-01 15:32:32.304482	2026-01-02 15:32:32.303815	\N	2026-01-01 15:32:32.304482	[::1]:51677	PostmanRuntime/7.51.0
11	1	d7d88a0b-44ac-40ac-bac7-007167059aa2	2026-01-01 15:33:12.282859	2026-01-02 15:33:12.282074	\N	2026-01-01 15:33:12.282859	[::1]:51677	PostmanRuntime/7.51.0
12	1	2fa48199-e489-436e-9e21-d11e157f3557	2026-01-01 15:45:18.887815	2026-01-02 15:45:18.886554	\N	2026-01-01 15:45:18.887815	[::1]:51792	PostmanRuntime/7.51.0
13	3	7f143cee-7db3-44d4-9ce0-2053b6799822	2026-01-01 15:48:11.850153	2026-01-02 15:48:11.849446	\N	2026-01-01 15:48:11.850153	[::1]:52102	PostmanRuntime/7.51.0
14	2	152a6a06-e081-4521-a4a8-b2cb41caa752	2026-01-03 16:10:05.737422	2026-01-04 16:10:05.728582	\N	2026-01-03 16:10:05.737422	[::1]:54129	PostmanRuntime/7.51.0
15	1	0938dcbe-bbeb-4512-beee-430ba38abd74	2026-01-03 16:10:44.205656	2026-01-04 16:10:44.204682	\N	2026-01-03 16:10:44.205656	[::1]:54129	PostmanRuntime/7.51.0
16	2	40c9c633-7d3b-4c7f-8528-b77373816235	2026-01-03 16:33:36.459682	2026-01-04 16:33:36.457362	\N	2026-01-03 16:33:36.459682	[::1]:54354	PostmanRuntime/7.51.0
17	2	3c4d5395-60b2-43e5-9358-f8ba170e7d90	2026-01-03 16:39:34.636222	2026-01-04 16:39:34.634866	\N	2026-01-03 16:39:34.636222	[::1]:54475	PostmanRuntime/7.51.0
18	5	43521023-6a85-4c76-9ae2-dfcd10775e7d	2026-01-03 16:55:13.693683	2026-01-04 16:55:13.682578	\N	2026-01-03 16:55:13.693683	[::1]:54628	PostmanRuntime/7.51.0
19	1	e090bd99-c5d7-49d0-8f72-6ad32c3a3049	2026-01-03 21:14:06.098205	2026-01-04 21:14:06.091031	\N	2026-01-03 21:14:06.098205	[::1]:57203	PostmanRuntime/7.51.0
20	5	212f5201-1457-4c55-ac1b-7f94a9aa803f	2026-01-03 21:41:10.878809	2026-01-04 21:41:10.868201	\N	2026-01-03 21:41:10.878809	[::1]:57525	PostmanRuntime/7.51.0
21	1	78a57d61-e9d8-4c9f-9dda-9ac1ea9b0986	2026-01-03 22:27:45.724417	2026-01-04 22:27:45.720369	\N	2026-01-03 22:27:45.724417	[::1]:57915	PostmanRuntime/7.51.0
22	2	9c9bc59e-ac1c-4616-a433-51b85262905d	2026-01-03 22:58:02.338119	2026-01-04 22:58:02.33687	\N	2026-01-03 22:58:02.338119	[::1]:58239	PostmanRuntime/7.51.0
23	1	818a6462-9530-4918-94b6-94100fcff21d	2026-01-04 00:36:46.376596	2026-01-05 00:36:46.367264	\N	2026-01-04 00:36:46.376596	[::1]:62466	PostmanRuntime/7.51.0
24	1	005f5a6b-0b4a-4b1e-8fa7-2dd4f75cb4af	2026-01-04 18:59:28.8521	2026-01-05 18:59:28.845489	\N	2026-01-04 18:59:28.8521	[::1]:53093	PostmanRuntime/7.51.0
25	1	594e399d-276b-444e-b404-9e1afe652c35	2026-01-04 20:19:51.311669	2026-01-05 20:19:51.310741	\N	2026-01-04 20:19:51.311669	[::1]:55768	PostmanRuntime/7.51.0
26	2	7b8c50db-a5ff-4644-b5c9-c0ca50a1a818	2026-01-05 13:29:48.5372	2026-01-06 13:29:48.534327	\N	2026-01-05 13:29:48.5372	[::1]:49996	PostmanRuntime/7.51.0
27	2	d50c5a59-ea4f-4804-9a8d-c2caa5b2d4d4	2026-01-08 00:25:43.979818	2026-01-09 00:25:43.97016	\N	2026-01-08 00:25:43.979818	[::1]:58836	PostmanRuntime/7.51.0
28	1	247ce304-ac09-4ab7-81ef-a4b58893ac59	2026-01-08 00:46:16.277734	2026-01-09 00:46:16.269776	\N	2026-01-08 00:46:16.277734	[::1]:59049	PostmanRuntime/7.51.0
29	1	5907ba49-be81-414b-bf23-2d7ce0f0f308	2026-01-08 11:52:49.571711	2026-01-09 11:52:49.565218	\N	2026-01-08 11:52:49.571711	[::1]:59547	PostmanRuntime/7.51.0
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password_hash, full_name, role, is_active, created_at, updated_at) FROM stdin;
2	johnstaff	john.staff@company.com	$2a$10$ZvxWDmKSX6dLxqZkqogTdeHhRYzezf53mz7GBO1fmsjwmfii7HUAa	John Staff	staff	t	2025-12-30 17:09:58.655162	2025-12-30 17:09:58.655162
1	superadmin	superadmin@inventory.com	$2a$10$pL9g8SnW/L.gUzdJ31XzD.F.dXU.m/681noh9vBfYItae8TXXeh1.	Super Administrator	super_admin	t	2025-12-30 16:24:13.626667	2025-12-30 17:59:14.948955
3	userfromsuperadmin	userfromsuperadmin@company.com	$2a$10$43KIdl0CWvxa.0O/upKfmukJ4vIUSLpVwm8GUt9VwU1ZUHD1S5epy	User From Super Admin	staff	t	2025-12-31 14:01:22.047577	2025-12-31 14:01:22.047577
5	admin	admin@inventory.com	$2a$10$AGdlOAt9bTkZdRXWT0FQq.F1U9zkJRSQzJMzCu7aqYzW9Tw6i76ti	Admin Administrator	admin	t	2026-01-03 16:54:25.633441	2026-01-03 16:54:25.633441
6	johnstaffupdated	john.updated@mail.com	$2a$10$HotBMbF292E3uOqw/SDDIOsrXV2/ItQ24K58fuRCdHfAyf3hd7VTe	John Doe Updated	staff	t	2026-01-08 11:54:55.32678	2026-01-08 11:56:44.859901
7	johndelete	johndelete@company.com	$2a$10$onMCm2cSnx.WeJdMkP0DreFitLxouyienZZRaQAri3RRkVya.Lq.2	Staff untuk testing delete	staff	f	2026-01-08 11:57:54.021072	2026-01-08 11:59:08.543443
\.


--
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.warehouses (id, code, name, address, city, province, postal_code, phone, is_active, created_by, created_at, updated_at) FROM stdin;
1	WH-001	Gudang Utama update	Jl. Industri No. 10	Jakarta	DKI Jakarta	12950	0217891234	t	1	2026-01-03 22:33:27.147194	2026-01-04 00:37:08.07825
2	BGR-001	Gudang Bogor	Jl. Surya Kencana No. 10	Bogor	Jawa Barat	13120	0217894321	t	1	2026-01-04 00:44:07.940452	2026-01-04 00:44:07.940452
3	DPR-001	Gudang DPR	Jl. Kencana No. 10	Neraka	Neraka paling bawah	13121	0217891111	f	1	2026-01-04 00:45:55.496112	2026-01-04 00:46:08.401677
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 4, true);


--
-- Name: items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.items_id_seq', 4, true);


--
-- Name: racks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.racks_id_seq', 4, true);


--
-- Name: sale_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sale_items_id_seq', 6, true);


--
-- Name: sales_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sales_id_seq', 8, true);


--
-- Name: sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sessions_id_seq', 29, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 7, true);


--
-- Name: warehouses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.warehouses_id_seq', 3, true);


--
-- Name: categories categories_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_code_key UNIQUE (code);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- Name: items items_sku_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_sku_key UNIQUE (sku);


--
-- Name: racks racks_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_code_key UNIQUE (code);


--
-- Name: racks racks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_pkey PRIMARY KEY (id);


--
-- Name: sale_items sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_pkey PRIMARY KEY (id);


--
-- Name: sales sales_invoice_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_invoice_number_key UNIQUE (invoice_number);


--
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_token_key UNIQUE (token);


--
-- Name: categories unique_categories_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT unique_categories_name UNIQUE (name);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: warehouses warehouses_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_code_key UNIQUE (code);


--
-- Name: warehouses warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);


--
-- Name: idx_categories_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_categories_code ON public.categories USING btree (code);


--
-- Name: idx_categories_is_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_categories_is_active ON public.categories USING btree (is_active);


--
-- Name: idx_items_category_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_category_id ON public.items USING btree (category_id);


--
-- Name: idx_items_is_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_is_active ON public.items USING btree (is_active);


--
-- Name: idx_items_minimum_stock; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_minimum_stock ON public.items USING btree (stock, minimum_stock);


--
-- Name: idx_items_rack_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_rack_id ON public.items USING btree (rack_id);


--
-- Name: idx_items_sku; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_sku ON public.items USING btree (sku);


--
-- Name: idx_items_stock; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_items_stock ON public.items USING btree (stock);


--
-- Name: idx_racks_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_racks_code ON public.racks USING btree (code);


--
-- Name: idx_racks_is_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_racks_is_active ON public.racks USING btree (is_active);


--
-- Name: idx_racks_warehouse_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_racks_warehouse_id ON public.racks USING btree (warehouse_id);


--
-- Name: idx_sale_items_item_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sale_items_item_id ON public.sale_items USING btree (item_id);


--
-- Name: idx_sale_items_sale_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sale_items_sale_id ON public.sale_items USING btree (sale_id);


--
-- Name: idx_sales_created_by; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_created_by ON public.sales USING btree (created_by);


--
-- Name: idx_sales_invoice_number; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_invoice_number ON public.sales USING btree (invoice_number);


--
-- Name: idx_sales_payment_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_payment_status ON public.sales USING btree (payment_status);


--
-- Name: idx_sales_sale_date; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_sale_date ON public.sales USING btree (sale_date);


--
-- Name: idx_sessions_expired_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_expired_at ON public.sessions USING btree (expired_at);


--
-- Name: idx_sessions_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_token ON public.sessions USING btree (token);


--
-- Name: idx_sessions_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_user_id ON public.sessions USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: idx_warehouses_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_warehouses_code ON public.warehouses USING btree (code);


--
-- Name: idx_warehouses_is_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_warehouses_is_active ON public.warehouses USING btree (is_active);


--
-- Name: categories update_categories_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: items update_items_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON public.items FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: racks update_racks_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_racks_updated_at BEFORE UPDATE ON public.racks FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: sales update_sales_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_sales_updated_at BEFORE UPDATE ON public.sales FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: warehouses update_warehouses_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_warehouses_updated_at BEFORE UPDATE ON public.warehouses FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: categories categories_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: items items_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE RESTRICT;


--
-- Name: items items_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: items items_rack_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_rack_id_fkey FOREIGN KEY (rack_id) REFERENCES public.racks(id) ON DELETE SET NULL;


--
-- Name: racks racks_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: racks racks_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE CASCADE;


--
-- Name: sale_items sale_items_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id) ON DELETE RESTRICT;


--
-- Name: sale_items sale_items_sale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_items
    ADD CONSTRAINT sale_items_sale_id_fkey FOREIGN KEY (sale_id) REFERENCES public.sales(id) ON DELETE CASCADE;


--
-- Name: sales sales_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: warehouses warehouses_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

