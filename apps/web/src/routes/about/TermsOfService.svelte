<script lang="ts">
	import { marked } from 'marked';
	import DOMPurify from 'dompurify';

	let html = $state('');

	let content = `
# PODMÍNKY UŽÍVÁNÍ PLATFORMY

## 1. Úvodní ustanovení

Tyto Podmínky užívání (dále jen „Podmínky“) upravují práva a povinnosti uživatelů a provozovatele při užívání vzdělávací platformy.

Užíváním platformy vyjadřuje uživatel souhlas s těmito Podmínkami a zavazuje se je dodržovat.

Provozovatelem Platformy je:

Název subjektu: __________________________ \n
IČO: __________________________ \n
Sídlo: __________________________ \n
Kontaktní e-mail: __________________________ \n

(dále jen „Provozovatel“)


## 2. Vznik uživatelského účtu

 ### 2.1 Oprávnění k registraci

Uživatelský účet může založit jakákoli fyzická osoba, která poskytne požadované registrační údaje a souhlasí s těmito Podmínkami.

Registrací vzniká mezi uživatelem a Provozovatelem smluvní vztah spočívající v poskytování vzdělávacích služeb prostřednictvím platformy.


## 3. Uživatelský účet a bezpečnost

 ### 3.1 Odpovědnost za přístupové údaje

Uživatel je plně odpovědný za ochranu svých přístupových údajů, zejména hesla.

- Uživatel je povinen uchovávat heslo v tajnosti.
- Uživatel nese odpovědnost za veškeré aktivity provedené prostřednictvím svého účtu.

Reset hesla automatizovaným způsobem není umožněn. V případě ztráty přístupu může uživatel kontaktovat správce systému, který situaci individuálně posoudí.

 ### 3.2 Zákaz sdílení přístupu

Uživatel není oprávněn sdílet své přístupové údaje s třetími osobami ani umožnit třetím osobám přístup ke svému účtu.

Porušení tohoto ustanovení může vést k omezení nebo zrušení uživatelského účtu.


## 4. Přístup k obsahu a role uživatelů

Platforma rozlišuje zejména následující role:

- Student
- Lektor
- případně Správce systému

 ### 4.1 Oprávnění lektora

Lektor je oprávněn:

- vytvářet a spravovat kurzy,
- nahrávat a upravovat studijní materiály,
- zamykat a odemykat jednotlivé bloky nebo obsah kurzu,
- vyhodnocovat testy a evidovat výsledky studentů.

 ### 4.2 Oprávnění studenta

Student má přístup pouze k těm kurzům, blokům a materiálům, které mu byly lektorem zpřístupněny.

Student nemá oprávnění zasahovat do struktury kurzu, upravovat obsah ani obcházet nastavená omezení přístupu.


## 5. Zrušení a omezení účtu

 ### 5.1 Zrušení účtu uživatelem

Uživatel může požádat o zrušení svého účtu, zejména z důvodu ochrany osobních údajů dle platné legislativy (GDPR).

Po zrušení účtu dojde k odstranění osobních údajů v souladu se Zásadami ochrany osobních údajů, pokud neexistuje zákonný důvod pro jejich další uchování.

 ### 5.2 Zrušení nebo omezení účtu Provozovatelem

Provozovatel si vyhrazuje právo omezit nebo zrušit uživatelský účet v případě, že:

- uživatel poruší tyto Podmínky,
- dojde ke zneužití účtu,
- uživatel neoprávněně sdílí přístupové údaje,
- uživatel porušuje právní předpisy.


## 6. Odpovědnost za obsah materiálů

Odpovědnost za obsah studijních materiálů, kurzů a testů nese primárně lektor, který daný obsah vytvořil nebo nahrál.

Provozovatel odpovídá za technický provoz platformy, nikoli za odbornou správnost nebo aktuálnost výukového obsahu, pokud jej sám nevytvořil.

Uživatelé jsou povinni zdržet se nahrávání obsahu, který:

- porušuje právní předpisy,
- zasahuje do práv třetích osob,
- explicitní obsah.


## 7. Ochrana autorských práv

Veškerý obsah zpřístupněný na platformě, zejména:

- studijní materiály,
- texty,
- prezentace,
- videa,
- testy,
- grafické prvky,

je chráněn autorským právem.

Uživatel není oprávněn obsah:

- kopírovat,
- reprodukovat,
- šířit,
- zpřístupňovat třetím osobám,
- využívat ke komerčním účelům,

bez předchozího souhlasu oprávněné osoby.

Autorskoprávní ochrana se řídí zákonem č. 121/2000 Sb., autorský zákon.

Porušení autorských práv může být řešeno v souladu s platnými právními předpisy.


## 8. Omezení odpovědnosti

Provozovatel nenese odpovědnost za:

- dočasnou nedostupnost platformy,
- ztrátu dat způsobenou vyšší mocí,
- škody vzniklé neoprávněným přístupem třetí osoby v důsledku nedostatečné ochrany přístupových údajů ze strany uživatele.

Provozovatel se zavazuje vyvíjet přiměřené úsilí k zajištění bezpečnosti a funkčnosti systému.


## 9. Závěrečná ustanovení

Tyto Podmínky mohou být přiměřeně měněny nebo aktualizovány, zejména v návaznosti na změny právních předpisů nebo rozvoj funkcionalit platformy.

Aktuální verze Podmínek je vždy zveřejněna na webových stránkách platformy.

Tyto Podmínky nabývají účinnosti dnem jejich zveřejnění: 02.03.2026.
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
