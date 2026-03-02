<script lang="ts">
	import { marked } from 'marked';
	import DOMPurify from 'dompurify';

	let html = $state('');

	let content = `
# ZÁSADY OCHRANY OSOBNÍCH ÚDAJŮ

## 1. Úvodní ustanovení

Tyto Zásady ochrany osobních údajů upravují způsob zpracování osobních údajů uživatelů vzdělávací platformy.

Provozovatel platformy při zpracování osobních údajů postupuje v souladu s platnými právními předpisy, zejména s:

- nařízením Evropského parlamentu a Rady (EU) 2016/679, obecné nařízení o ochraně osobních údajů (General Data Protection Regulation, dále jen „GDPR“),
- zákonem č. 110/2019 Sb., o zpracování osobních údajů.


## 2. Správce osobních údajů

Správcem osobních údajů je:

Název subjektu: __________________________ \n
IČO: __________________________ \n
Sídlo: __________________________ \n
Kontaktní e-mail: __________________________ \n

(dále jen „Správce“)

Správce určuje účely a prostředky zpracování osobních údajů a odpovídá za jejich řádné a zákonné zpracování.


## 3. Rozsah zpracovávaných osobních údajů

Správce zpracovává osobní údaje, které mu uživatel poskytne při registraci, užívání platformy nebo v souvislosti s poskytováním vzdělávacích služeb.

Zpracovávány jsou zejména tyto údaje:

- jméno a příjmení,
- e-mailová adresa,
- údaje o uživatelském účtu,
- přístupové údaje,
- výsledky testů a studijní pokrok,
- studijní materiály nahrané v rámci kurzů.

Osobní údaje jsou zpracovávány v elektronické podobě prostřednictvím zabezpečených informačních systémů.


## 4. Účel zpracování osobních údajů

Osobní údaje jsou zpracovávány za účelem:

- zajištění registrace a správy uživatelských účtů,
- poskytování vzdělávacích služeb prostřednictvím platformy,
- umožnění přístupu ke studijním materiálům,
- vyhodnocování testů a evidence studijních výsledků,
- oznámení lektorů běžným uživatelům.

Právním základem zpracování je plnění smlouvy o poskytování služeb mezi Správcem a uživatelem ve smyslu čl. 6 odst. 1 písm. b) GDPR.


## 5. Doba uchovávání osobních údajů

Osobní údaje jsou uchovávány po dobu trvání uživatelského účtu.

- Údaje související s uživatelským účtem (včetně přístupových údajů) jsou uchovávány do doby zrušení účtu uživatelem nebo Správcem.
- Studijní materiály a výsledky testů jsou uchovávány po dobu jejich dostupnosti v rámci příslušného kurzu, případně do jejich odstranění lektorem.
- Po zrušení účtu jsou osobní údaje bez zbytečného odkladu vymazány, pokud neexistuje zákonný důvod pro jejich další uchování.


## 6. Přístup k osobním údajům

K osobním údajům mají přístup pouze oprávněné osoby v rámci Správce a osoby podílející se na provozu platformy.

V rámci systému přístupových práv:

- Lektor má přístup k údajům vztahujícím se ke kurzům, které spravuje, včetně studijních materiálů a výsledků testů studentů.
- Student má přístup pouze ke kurzům, materiálům a výsledkům, které mu byly lektorem zpřístupněny.

Osobní údaje nejsou předávány třetím osobám, s výjimkou případů, kdy to vyžaduje právní předpis.


## 7. Zabezpečení osobních údajů

Správce přijal přiměřená technická a organizační opatření k ochraně osobních údajů proti neoprávněnému přístupu, ztrátě, zneužití nebo poškození.

Přístupové údaje (hesla) jsou ukládány v zabezpečené podobě.


## 8. Práva subjektů údajů

Uživatel jako subjekt údajů má v souladu s GDPR právo:

- na přístup ke svým osobním údajům,
- na opravu nepřesných nebo neúplných údajů,
- na výmaz osobních údajů,
- na omezení zpracování,
- vznést námitku proti zpracování,
- na přenositelnost údajů,
- podat stížnost u dozorového úřadu.

Dozorovým úřadem v České republice je Úřad pro ochranu osobních údajů (ÚOOÚ).

Žádosti týkající se uplatnění práv lze zasílat na kontaktní e-mail Správce uvedený výše.


## 9. Cookies

Platforma používá základní (technické) cookies nezbytné pro její správnou funkčnost, zejména pro:

- přihlášení uživatele,
- správu relace,
- zabezpečení přístupu k účtu.

Tyto cookies nevyžadují souhlas uživatele, jelikož jsou nezbytné pro fungování služby.


## 10. Závěrečná ustanovení

Tyto Zásady ochrany osobních údajů mohou být Správcem přiměřeně aktualizovány v návaznosti na změny právních předpisů nebo změny fungování platformy.

Aktuální verze je vždy zveřejněna na webových stránkách platformy.

Verze 02.03.2026

`;

	import { onMount } from 'svelte';
	onMount(async () => {
		// This ONLY runs in the browser, solving hydration issues
		const raw = await marked.parse(content);
		html = DOMPurify.sanitize(raw);
	});

	// const rawHtml = marked.parse(content) as string;
	// const html = DOMPurify.sanitize(rawHtml);
</script>

<div class="relative mx-auto my-12 w-full max-w-4xl">
	<div class="absolute inset-0 translate-x-3 translate-y-3 rounded-2xl bg-s-black"></div>

	<div
		class="relative overflow-hidden rounded-2xl border-4 border-s-black bg-white p-8 shadow-none md:p-12"
	>
		<div class="mb-10 text-center">
			<div
				class="mb-4 inline-block -rotate-1 border-4 border-s-black bg-p-green px-6 py-2 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
			>
				<span class="text-xl font-black tracking-tighter text-s-black uppercase"
					>Official Document</span
				>
			</div>
			<div class="mx-auto h-2 w-32 bg-p-blue"></div>
		</div>

		<div
			class="legal-content
            [&_h1]:mb-8 [&_h1]:text-5xl [&_h1]:font-black [&_h1]:tracking-tighter [&_h1]:text-s-black [&_h1]:uppercase
            [&_h2]:mt-12 [&_h2]:mb-4 [&_h2]:inline-block [&_h2]:border-2 [&_h2]:border-s-black [&_h2]:bg-p-blue [&_h2]:px-3 [&_h2]:py-1 [&_h2]:text-2xl [&_h2]:font-black [&_h2]:text-white [&_h2]:uppercase
            [&_h3]:mt-8 [&_h3]:mb-2 [&_h3]:text-xl [&_h3]:font-black [&_h3]:text-s-black [&_h3]:uppercase
            [&_li]:mb-3 [&_li]:font-bold [&_li]:before:mr-2 [&_li]:before:content-['•'] [&_p]:mb-6
            [&_p]:text-lg [&_p]:leading-relaxed [&_p]:font-bold
            [&_p]:text-s-black/80 [&_strong]:font-black [&_strong]:text-p-blue [&_ul]:mb-6
            [&_ul]:list-none [&_ul]:pl-0"
		>
			{@html html}
		</div>

		<div class="mt-12 border-t-4 border-dashed border-s-black/20 pt-8 text-center">
			<p class="text-sm font-black tracking-widest text-s-black/40 uppercase">End of Document</p>
		</div>
	</div>
</div>
